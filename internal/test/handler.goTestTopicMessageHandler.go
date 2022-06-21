package standardtest

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

func (h *GoTestTopicMessageHandler) ProcessMessage(message *nsq.Message) error {
	log.Printf("Message on %s (%s): %v\n", message.Topic, message.NSQDAddress, string(message.Body))

	return message.ForwardUnhandledMessageHandler()
}
