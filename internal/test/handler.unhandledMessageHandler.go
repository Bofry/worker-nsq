package test

import (
	nsq "github.com/Bofry/worker-nsq"
)

type UnhandledMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *UnhandledMessageHandler) Init() {
	defaultLogger.Printf("UnhandledMessageHandler.Init()")
}

func (h *UnhandledMessageHandler) ProcessMessage(message *nsq.Message) error {
	defaultLogger.Printf("Unhandled Message on %s (%s): %v\n", message.Topic, message.NSQDAddress, string(message.Body))

	return nil
}
