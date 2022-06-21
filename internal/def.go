package internal

import (
	"log"
	"os"
	"reflect"

	nsq "github.com/Bofry/lib-nsq"
)

const (
	LOGGER_PREFIX string = "[worker-nsq] "
)

var (
	NsqWorkerServiceInstance = new(NsqWorkerService)

	typeOfHost = reflect.TypeOf(NsqWorker{})

	logger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	Config            = nsq.Config
	Message           = nsq.Message
	MessageHandleProc = nsq.MessageHandleProc

	MessageHandler interface {
		ProcessMessage(message *Message) error
	}
)
