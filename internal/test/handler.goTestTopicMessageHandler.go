package test

import (
	"log"

	nsq "github.com/Bofry/worker-nsq"
)

type GoTestTopicMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *GoTestTopicMessageHandler) Init() {
	log.Printf("GoTestTopicMessageHandler.Init()")
}

func (h *GoTestTopicMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	log.Printf("Message on %s (%s): [%s] %v\n", message.Topic, message.NSQDAddress, message.ID, string(message.Body))

	return ctx.ForwardUnhandledMessage(message)
}
