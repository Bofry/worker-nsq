package test

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/config"
	nsq "github.com/Bofry/worker-nsq"
)

func TestStarter(t *testing.T) {
	err := setupTestStarter()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := teardownTestStarter()
		if err != nil {
			t.Fatal(err)
		}
	}()

	app := App{}
	starter := nsq.Startup(&app).
		Middlewares(
			nsq.UseTopicGateway(&TopicGateway{}),
			// nsq.UseErrorHandler(func(err error) (disposed bool) {
			// 	t.Logf("catch err: %v", err)
			// 	return false
			// }),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
		var expectedNsqAddress string = os.Getenv("NSQ_ADDRESS")
		if !reflect.DeepEqual(conf.NsqAddress, expectedNsqAddress) {
			t.Errorf("assert 'Config.RedisAddress':: expected '%v', got '%v'", expectedNsqAddress, conf.NsqAddress)
		}
	}
}

func setupTestStarter() error {
	p, err := nsq.NewForwarder(&nsq.ProducerOption{
		Address: strings.Split(os.Getenv("NSQD_SERVERS"), ","),
		Config:  nsq.NewConfig(),
	})
	if err != nil {
		return err
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

	return nil
}

func teardownTestStarter() error {
	return nil
}
