package test

import (
	"fmt"
	"reflect"

	nsq "github.com/Bofry/worker-nsq"
)

var (
	_ nsq.MessageHandler        = new(GoTestTopicMessageHandler)
	_ nsq.MessageObserverAffair = new(GoTestTopicMessageHandler)
)

type BlackholdMessageHandler struct {
}

func (h *BlackholdMessageHandler) Init() {
	fmt.Println("BlackholdMessageHandler.Init()")
}

func (h *BlackholdMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	// disable recording log
	ctx.RecordingLog(false)

	return nil
}

// MessageObserverTypes implements internal.MessageObserverAffair.
func (*BlackholdMessageHandler) MessageObserverTypes() []reflect.Type {
	return []reflect.Type{}
}
