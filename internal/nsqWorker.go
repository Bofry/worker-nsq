package internal

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/Bofry/host"
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var _ host.Host = new(NsqWorker)

type NsqWorker struct {
	NsqAddress         string // nsqd:127.0.0.1:4150,127.0.0.2:4150 -or- nsqlookupd:127.0.0.1:4161,127.0.0.2:4161
	Channel            string
	HandlerConcurrency int
	Config             *Config

	consumer *nsq.Consumer

	logger *log.Logger

	messageDispatcher *MessageDispatcher
	messageManager    interface{}

	messageHandleService *MessageHandleService
	messageTracerService *MessageTracerService

	onErrorEventHandler host.HostOnErrorEventHandler

	wg          sync.WaitGroup
	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (w *NsqWorker) Start(ctx context.Context) {
	if w.disposed {
		NsqWorkerLogger.Panic("the Worker has been disposed")
	}
	if !w.initialized {
		NsqWorkerLogger.Panic("the Worker havn't be initialized yet")
	}
	if w.running {
		return
	}

	var err error
	w.mutex.Lock()
	defer func() {
		if err != nil {
			w.running = false
			w.disposed = true
		}
		w.mutex.Unlock()
	}()

	w.running = true
	w.messageDispatcher.start(ctx)

	var (
		topics = w.messageDispatcher.Topics()
	)

	NsqWorkerLogger.Printf("channel [%s] topics [%s] on address %s\n",
		w.Channel,
		strings.Join(topics, ","),
		w.NsqAddress)

	if len(topics) > 0 {
		c := w.consumer
		err := c.Subscribe(topics)
		if err != nil {
			NsqWorkerLogger.Panic(err)
		}
	}
}

func (w *NsqWorker) Stop(ctx context.Context) error {
	NsqWorkerLogger.Printf("%% Stopping\n")

	w.mutex.Lock()
	defer func() {
		w.running = false
		w.disposed = true
		w.mutex.Unlock()

		w.messageDispatcher.stop(ctx)

		NsqWorkerLogger.Printf("%% Stopped\n")
	}()

	w.consumer.Close()
	w.wg.Wait()
	return nil
}

func (w *NsqWorker) Logger() *log.Logger {
	return w.logger
}

func (w *NsqWorker) alloc() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.messageHandleService = NewMessageHandleService()
	w.messageTracerService = NewMessageTracerService()

	w.messageDispatcher = &MessageDispatcher{
		MessageHandleService: w.messageHandleService,
		MessageTracerService: w.messageTracerService,
		Router:               make(Router),
		OnHostErrorProc:      w.onHostError,
	}
}

func (w *NsqWorker) init() {
	if w.initialized {
		return
	}

	w.mutex.Lock()
	defer func() {
		w.initialized = true
		w.mutex.Unlock()
	}()

	w.messageTracerService.init(w.messageManager)
	w.messageDispatcher.init()
	w.configConsumer()
}

func (w *NsqWorker) configConsumer() {
	instance := &nsq.Consumer{
		NsqAddress:         w.NsqAddress,
		Channel:            w.Channel,
		HandlerConcurrency: w.HandlerConcurrency,
		Config:             w.Config,
		MessageHandler:     w.receiveMessage,
		Logger:             w.logger,
	}

	w.consumer = instance
}

func (w *NsqWorker) receiveMessage(message *Message) error {
	ctx := &Context{
		Channel:                 w.Channel,
		logger:                  w.logger,
		unhandledMessageHandler: nil, // be determined by MessageDispatcher
	}
	return w.messageDispatcher.ProcessMessage(ctx, message)
}

func (w *NsqWorker) onHostError(err error) (disposed bool) {
	if w.onErrorEventHandler != nil {
		return w.onErrorEventHandler.OnError(err)
	}
	return false
}

func (w *NsqWorker) setTextMapPropagator(propagator propagation.TextMapPropagator) {
	w.messageTracerService.TextMapPropagator = propagator
}

func (w *NsqWorker) setTracerProvider(provider *trace.SeverityTracerProvider) {
	w.messageTracerService.TracerProvider = provider
}

func (w *NsqWorker) setLogger(l *log.Logger) {
	w.logger = l
}
