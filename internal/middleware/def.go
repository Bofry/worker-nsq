package middleware

import (
	"reflect"

	"github.com/Bofry/worker-nsq/internal"
)

const (
	UNHANDLED_MESSAGE_HANDLER_TOPIC_SYMBOL string = "?"

	TAG_TOPIC = "topic"
)

var (
	typeOfHost           = reflect.TypeOf(internal.NsqWorker{})
	typeOfMessageHandler = reflect.TypeOf((*internal.MessageHandler)(nil)).Elem()
)

type (
	ConfigureUnhandledMessageHandleProc func(handler internal.MessageHandler)
)
