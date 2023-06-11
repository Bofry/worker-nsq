package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/config"
	nsq "github.com/Bofry/worker-nsq"
	"github.com/joho/godotenv"
)

var (
	__ENV_FILE        = ".env"
	__ENV_FILE_SAMPLE = ".env.sample"

	__CONFIG_YAML_FILE        = "config.yaml"
	__CONFIG_YAML_FILE_SAMPLE = "config.yaml.sample"
)

type MessageManager struct {
	GoTest2Topic *GoTestTopicMessageHandler `topic:"gotest2Topic"`
	GoTestTopic  *GoTestTopicMessageHandler `topic:"gotestTopic"`
	Unhandled    *UnhandledMessageHandler   `topic:"?"`
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func TestMain(m *testing.M) {
	var err error

	_, err = os.Stat(__CONFIG_YAML_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__CONFIG_YAML_FILE_SAMPLE, __CONFIG_YAML_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	_, err = os.Stat(__ENV_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__ENV_FILE_SAMPLE, __ENV_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	godotenv.Load(__ENV_FILE)
	{
		p, err := nsq.NewForwarder(&nsq.ProducerConfig{
			Address: strings.Split(os.Getenv("TEST_NSQD_SERVERS"), ","),
			Config:  nsq.NewConfig(),
			Logger:  defaultLogger,
		})
		if err != nil {
			panic(err)
		}

		{
			topic := "gotestTopic"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.Write(topic, []byte(word))
			}
		}

		{
			topic := "gotest2Topic"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.Write(topic, []byte(word))
			}
		}

		{
			topic := "unknownTopic"
			for _, word := range []string{"unknown"} {
				p.Write(topic, []byte(word))
			}
		}

		p.Close()
	}

	m.Run()
}

func TestStartup(t *testing.T) {
	app := App{}
	starter := nsq.Startup(&app).
		Middlewares(
			nsq.UseMessageManager(&MessageManager{}),
			nsq.UseErrorHandler(func(ctx *nsq.Context, msg *nsq.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			nsq.UseTracing(false),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	// assert app.Config
	{
		conf := app.Config
		var expectedNsqAddress string = os.Getenv("TEST_NSQLOOKUPD_ADDRESS")
		if !reflect.DeepEqual(conf.NsqAddress, expectedNsqAddress) {
			t.Errorf("assert 'Config.NsqAddress':: expected '%v', got '%v'", expectedNsqAddress, conf.NsqAddress)
		}
	}
}

func TestStartup_UseTracing(t *testing.T) {
	var (
		testStartAt time.Time
	)

	app := App{}
	starter := nsq.Startup(&app).
		Middlewares(
			nsq.UseMessageManager(&MessageManager{}),
			nsq.UseErrorHandler(func(ctx *nsq.Context, msg *nsq.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			nsq.UseTracing(true),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	testStartAt = time.Now()

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		testEndAt := time.Now()

		// wait 2 seconds
		time.Sleep(2 * time.Second)

		var queryUrl = fmt.Sprintf(
			"%s?end=%d&limit=50&lookback=1h&&service=nsq-trace-demo&start=%d",
			app.Config.JaegerQueryUrl,
			testEndAt.UnixMicro(),
			testStartAt.UnixMicro())

		t.Log(queryUrl)
		req, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert query 'Jeager Query Url StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("%v", string(body))
		// parse content
		{
			var reply map[string]interface{}
			dec := json.NewDecoder(bytes.NewBuffer(body))
			dec.UseNumber()
			if err := dec.Decode(&reply); err != nil {
				t.Error(err)
			}

			data := reply["data"].([]interface{})
			if data == nil {
				t.Errorf("missing data section")
			}
			var expectedDataLength int = 14
			if expectedDataLength != len(data) {
				t.Errorf("assert 'Jaeger Query size of replies':: expected '%v', got '%v'", expectedDataLength, len(data))
			}
		}
	}
}

func TestStartup_UseLogging_And_UseTracing(t *testing.T) {
	var (
		loggingBuffer bytes.Buffer
	)

	app := App{}
	starter := nsq.Startup(&app).
		Middlewares(
			nsq.UseMessageManager(&MessageManager{}),
			nsq.UseErrorHandler(func(ctx *nsq.Context, msg *nsq.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			nsq.UseLogging(&MultiLoggerService{
				LoggingServices: []nsq.LoggingService{
					&LoggingService{},
					&BlackholeLoggerService{
						Buffer: &loggingBuffer,
					},
				},
			}),
			nsq.UseTracing(true),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		// test loggingBuffer
		var expectedLoggingBuffer string = strings.Join([]string{
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"BeforeProcessMessage()\n",
			"AfterProcessMessage()\n",
			"Flush()\n",
		}, "")
		if expectedLoggingBuffer != loggingBuffer.String() {
			t.Errorf("assert loggingBuffer:: expected '%v', got '%v'", expectedLoggingBuffer, loggingBuffer.String())
		}
	}
}
