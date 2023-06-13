package test

import (
	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.MessageHandler = new(GoTestTopicMessageHandler)

type GoTestTopicMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *GoTestTopicMessageHandler) Init() {
	defaultLogger.Printf("GoTestTopicMessageHandler.Init()")
}

func (h *GoTestTopicMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	defaultLogger.Printf("Message on %s (%s): [%s] %v\n", message.Topic, message.NSQDAddress, message.ID, string(message.Body))

	return ctx.ForwardUnhandledMessage(message)
}
