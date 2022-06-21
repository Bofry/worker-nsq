package standardtest

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
	}

	Host nsq.Worker

	Config struct {
		// redis
		NsqAddress            string `env:"*NSQ_ADDRESS"         yaml:"-"`
		NsqChannel            string `env:"-"                    yaml:"NsqChannel"`
		NsqHandlerConcurrency int    `env:"-"                    yaml:"NsqHandlerConcurrency"`
	}

	ServiceProvider struct {
		ResourceName string
	}

	TopicGateway struct {
		*GoTestTopicMessageHandler `topic:"gotestTopic"`
		*UnhandledMessageHandler   `topic:"?"`
	}
)

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
