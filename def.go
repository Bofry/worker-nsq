package nsq

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
	"github.com/Bofry/worker-nsq/internal/middleware"
)

const (
	// INVALID = internal.INVALID
	// UNSET   = internal.UNSET
	// PASS    = internal.PASS
	// FAIL    = internal.FAIL
	// ABORT   = internal.ABORT

	ContextStatusUnkonwn = internal.ContextStatusUnkonwn
	ContextStatusOK      = internal.ContextStatusOK
	ContextStatusFail    = internal.ContextStatusFail
)

type (
	ProducerConfig = nsq.ProducerConfig
	Producer       = nsq.Producer
	Forwarder      = nsq.Forwarder
	Message        = nsq.Message
	MessageContent = nsq.MessageContent
	Config         = nsq.Config

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
