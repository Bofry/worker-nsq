package nsq

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
	"github.com/Bofry/worker-nsq/internal/middleware"
)

const (
	StatusInvalid = internal.INVALID
	StatusUnset   = internal.UNSET
	StatusPass    = internal.PASS
	StatusFail    = internal.FAIL
	StatusAbort   = internal.ABORT
)

type (
	ProducerConfig = nsq.ProducerConfig
	Producer       = nsq.Producer
	Forwarder      = nsq.Forwarder
	Message        = nsq.Message
	MessageContent = nsq.MessageContent
	Config         = nsq.Config
	LogLevel       = nsq.LogLevel
	ConsumerOption = nsq.ConsumerOption

	EventEvidence  = middleware.EventEvidence
	LoggingService = middleware.LoggingService
	EventLog       = middleware.EventLog

	MessageObserver       = internal.MessageObserver
	MessageObserverAffair = internal.MessageObserverAffair

	MessageHandler = internal.MessageHandler
	Worker         = internal.NsqWorker
	Context        = internal.Context
	ReplyCode      = internal.ReplyCode

	ErrorHandler = internal.ErrorHandler
)

func NewForwarder(config *ProducerConfig) (*Forwarder, error) {
	return nsq.NewForwarder(config)
}

func NewConfig() *Config {
	return nsq.NewConfig()
}
