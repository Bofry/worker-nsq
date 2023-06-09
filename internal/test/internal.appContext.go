package test

import (
	"fmt"
	"time"

	nsq "github.com/Bofry/worker-nsq"
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
		// redis
		NsqAddress            string `env:"*NSQD_ADDRESS"        yaml:"-"`
		NsqChannel            string `env:"-"                    yaml:"NsqChannel"`
		NsqHandlerConcurrency int    `env:"-"                    yaml:"NsqHandlerConcurrency"`
	}

	ServiceProvider struct {
		ResourceName string
	}

	MessageManager struct {
		GoTest2Topic *GoTestTopicMessageHandler `topic:"gotest2Topic"`
		GoTestTopic  *GoTestTopicMessageHandler `topic:"gotestTopic"`
		Unhandled    *UnhandledMessageHandler   `topic:"?"`
	}
)

func (app *App) Init(conf *Config) {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
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
