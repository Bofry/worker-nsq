package nsq

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
)

type (
	ProducerOption = nsq.ProducerOption
	Forwarder      = nsq.Forwarder
	Message        = nsq.Message
	Config         = nsq.Config

	MessageHandler = internal.MessageHandler
	Worker         = internal.NsqWorker
)

func NewForwarder(opt *ProducerOption) (*Forwarder, error) {
	return nsq.NewForwarder(opt)
}

func NewConfig() *Config {
	return nsq.NewConfig()
}
