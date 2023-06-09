package nsq

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
	"github.com/Bofry/worker-nsq/internal/middleware"
)

type (
	ProducerConfig = nsq.ProducerConfig
	Producer       = nsq.Producer
	Forwarder      = nsq.Forwarder
	Message        = nsq.Message
	Config         = nsq.Config

	LoggingService = middleware.LoggingService

	MessageHandler = internal.MessageHandler
	Worker         = internal.NsqWorker
	Context        = internal.Context

	ErrorHandler = internal.ErrorHandler
)

func NewForwarder(config *ProducerConfig) (*Forwarder, error) {
	return nsq.NewForwarder(config)
}

func NewConfig() *Config {
	return nsq.NewConfig()
}
