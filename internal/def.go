package internal

import (
	"context"
	"log"
	"os"
	"reflect"

	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

const (
	LOGGER_PREFIX string = "[worker-nsq] "
)

var (
	typeOfHost               = reflect.TypeOf(NsqWorker{})
	defaultTracerProvider    = createNoopTracerProvider()
	defaultTextMapPropagator = createNoopTextMapPropagator()

	NsqWorkerModuleInstance = NsqWorkerModule{}

	NsqWorkerLogger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	Config  = nsq.Config
	Message = nsq.Message

	MessageHandleModule interface {
		CanSetSuccessor() bool
		SetSuccessor(successor MessageHandleModule)
		ProcessMessage(ctx *Context, message *Message, state ProcessingState, recover *Recover) error
		OnInitComplete()
		OnStart(ctx context.Context) error
		OnStop(ctx context.Context) error
	}

	MessageHandler interface {
		ProcessMessage(ctx *Context, message *Message) error
	}

	ErrorHandler func(ctx *Context, message *Message, err interface{})

	OnHostErrorHandler func(err error) (disposed bool)
)

func createNoopTracerProvider() *trace.SeverityTracerProvider {
	tp, err := trace.NoopProvider()
	if err != nil {
		NsqWorkerLogger.Fatalf("cannot create NoopProvider: %v", err)
	}
	return tp
}

func createNoopTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator()
}
