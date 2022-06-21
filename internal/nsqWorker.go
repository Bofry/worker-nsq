package internal

import (
	"context"
	"strings"
	"sync"

	"github.com/Bofry/host"
	nsq "github.com/Bofry/lib-nsq"
)

var _ host.Host = new(NsqWorker)

type NsqWorker struct {
	NsqAddress         string // nsqd:127.0.0.1:4150,127.0.0.2:4150 -or- nsqlookupd:127.0.0.1:4161,127.0.0.2:4161
	Channel            string
	HandlerConcurrency int
	Config             *Config

	consumer *nsq.Consumer

	dispatcher *NsqMessageDispatcher

	wg          sync.WaitGroup
	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (w *NsqWorker) Start(ctx context.Context) {
	if w.disposed {
		logger.Panic("the Worker has been disposed")
	}
	if !w.initialized {
		logger.Panic("the Worker havn't be initialized yet")
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

	var (
		topics = w.dispatcher.Topics()
	)

	logger.Printf("channel [%s] topics [%s] on address %s\n",
		w.Channel,
		strings.Join(topics, ","),
		w.NsqAddress)

	if len(topics) > 0 {
		c := w.consumer
		err := c.Subscribe(topics)
		if err != nil {
			logger.Panic(err)
		}
	}
}

func (w *NsqWorker) Stop(ctx context.Context) error {
	logger.Printf("%% Stopping\n")
	defer func() {
		logger.Printf("%% Stopped\n")
	}()

	w.consumer.Close()
	return nil
}

func (w *NsqWorker) preInit() {
	w.dispatcher = NewNsqMessageDispatcher()
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

	w.configConsumer()
}

func (w *NsqWorker) configConsumer() {
	instance := &nsq.Consumer{
		NsqAddress:              w.NsqAddress,
		Channel:                 w.Channel,
		HandlerConcurrency:      w.HandlerConcurrency,
		Config:                  w.Config,
		MessageHandler:          w.dispatcher.ProcessMessage,
		UnhandledMessageHandler: w.dispatcher.ProcessUnhandledMessage,
	}

	w.consumer = instance
}
