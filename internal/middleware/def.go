package middleware

import (
	"log"
	"reflect"

	"github.com/Bofry/worker-nsq/internal"
)

const (
	INVALID_MESSAGE_HANDLER_TOPIC_SYMBOL string = "?"

	TAG_TOPIC          = "topic"
	TAG_OPT_EXPAND_ENV = "@ExpandEnv"
	OPT_ON             = "on"
	OPT_OFF            = "off"
)

var (
	typeOfHost            = reflect.TypeOf(internal.NsqWorker{})
	typeOfMessageObserver = reflect.TypeOf((*internal.MessageObserver)(nil)).Elem()

	NoopLoggingServiceSingleton = NoopLoggingService{}
	NoopEventLogSingleton       = NoopEventLog(0)
)

type (
	LoggingService interface {
		CreateEventLog(ev EventEvidence) EventLog
		ConfigureLogger(l *log.Logger)
	}

	EventLog interface {
		OnError(message *internal.Message, err interface{}, stackTrace []byte)
		OnProcessMessage(message *internal.Message)
		OnProcessMessageComplete(message *internal.Message, reply internal.ReplyCode)
		Flush()
	}
)
