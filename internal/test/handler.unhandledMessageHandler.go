package test

import (
	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.MessageHandler = new(UnhandledMessageHandler)

type UnhandledMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *UnhandledMessageHandler) Init() {
	defaultLogger.Printf("UnhandledMessageHandler.Init()")
}

func (h *UnhandledMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	defaultLogger.Printf("Unhandled Message on %s (%s): %v\n", message.Topic, message.NSQDAddress, string(message.Body))

	return nil
}
