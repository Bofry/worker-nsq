package test

import (
	"context"

	"github.com/Bofry/trace"
	nsq "github.com/Bofry/worker-nsq"
	"github.com/Bofry/worker-nsq/tracing"
)

var _ nsq.MessageHandler = new(GoTestTopicMessageHandler)

type GoTestTopicMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *GoTestTopicMessageCounter
}

func (h *GoTestTopicMessageHandler) Init() {
	defaultLogger.Printf("GoTestTopicMessageHandler.Init()")

	h.counter = new(GoTestTopicMessageCounter)
}

func (h *GoTestTopicMessageHandler) ProcessMessage(ctx *nsq.Context, message *nsq.Message) error {
	defaultLogger.Printf("Message on %s (%s): [%s] %v\n", message.Topic, message.NSQDAddress, message.ID, string(message.Body))

	sp := trace.SpanFromContext(ctx)
	sp.Argv(string(message.Body))

	if message.Topic == "gotest2Topic" {
		h.doSomething(sp.Context())
		return ctx.ThrowInvalidMessageError(message)
	}
	h.counter.increase(sp.Context())
	return nil
}

func (h *GoTestTopicMessageHandler) doSomething(ctx context.Context) {
	tr := tracing.GetTracer(h)
	sp := tr.Start(ctx, "doSomething()")
	defer sp.End()
}

type GoTestTopicMessageCounter struct {
	count int
}

func (c *GoTestTopicMessageCounter) increase(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "increase()")
	defer sp.End()

	c.count++
	return c.count
}
