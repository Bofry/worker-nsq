package test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Bofry/trace"
	nsq "github.com/Bofry/worker-nsq"
	"go.opentelemetry.io/otel/propagation"
)

var (
	defaultLogger *log.Logger = log.New(log.Writer(), "[worker-nsq-test] ", log.LstdFlags|log.Lmsgprefix|log.LUTC)
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider

		Component       *MockComponent
		ComponentRunner *MockComponentRunner
	}

	Host nsq.Worker

	Config struct {
		// nsq
		NsqAddress            string `env:"*TEST_NSQLOOKUPD_ADDRESS"   yaml:"-"`
		NsqChannel            string `env:"-"                          yaml:"NsqChannel"`
		NsqHandlerConcurrency int    `env:"-"                          yaml:"NsqHandlerConcurrency"`

		// jaeger
		JaegerTraceUrl string `yaml:"jaegerTraceUrl"`
		JaegerQueryUrl string `yaml:"jaegerQueryUrl"`
	}

	ServiceProvider struct {
		ResourceName string
	}
)

func (app *App) Init(conf *Config) {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
}

func (app *App) OnInit() {
}

func (app *App) OnInitComplete() {
}

func (app *App) OnStart(ctx context.Context) {
}

func (app *App) OnStop(ctx context.Context) {
	{
		defaultLogger.Printf("stoping TracerProvider")
		tp := trace.GetTracerProvider()
		err := tp.Shutdown(ctx)
		if err != nil {
			defaultLogger.Printf("stoping TracerProvider error: %+v", err)
		}
	}
}

func (app *App) ConfigureLogger(l *log.Logger) {
	l.SetFlags(defaultLogger.Flags())
	l.SetOutput(defaultLogger.Writer())
}

func (app *App) Logger() *log.Logger {
	return defaultLogger
}

func (app *App) ConfigureTracerProvider() {
	if len(app.Config.JaegerTraceUrl) == 0 {
		tp, _ := trace.NoopProvider()
		trace.SetTracerProvider(tp)
		return
	}

	tp, err := trace.JaegerProvider(app.Config.JaegerTraceUrl,
		trace.ServiceName("nsq-trace-demo"),
		trace.Environment("go-bofry-worker-nsq-test"),
		trace.Pid(),
	)
	if err != nil {
		defaultLogger.Fatal(err)
	}

	trace.SetTracerProvider(tp)
}

func (app *App) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
}

func (app *App) ConfigureTextMapPropagator() {
	trace.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func (app *App) TextMapPropagator() propagation.TextMapPropagator {
	return trace.GetTextMapPropagator()
}

func (provider *ServiceProvider) Init(conf *Config) {
	fmt.Println("ServiceProvider.Init()")
	provider.ResourceName = "demo resource"
}

func (h *Host) Init(conf *Config) {
	config := nsq.NewConfig()
	{
		config.LookupdPollInterval = time.Second * 3
		config.DefaultRequeueDelay = 0
		config.MaxBackoffDuration = time.Millisecond * 50
		config.LowRdyIdleTimeout = time.Second * 1
		config.RDYRedistributeInterval = time.Millisecond * 20
	}

	h.NsqAddress = conf.NsqAddress
	h.Channel = conf.NsqChannel
	h.HandlerConcurrency = conf.NsqHandlerConcurrency
	h.Config = config
}

func (h *Host) OnError(err error) (disposed bool) {
	return false
}
