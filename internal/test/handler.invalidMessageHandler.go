package test

import (
	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.MessageHandler = new(InvalidMessageHandler)

type InvalidMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *InvalidMessageHandler) Init() {
	defaultLogger.Printf("InvalidMessageHandler.Init()")
}

func (h *InvalidMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	defaultLogger.Printf("Invalid Message on %s (%s): %v\n", message.Topic, message.NSQDAddress, string(message.Body))

	return nil
}
