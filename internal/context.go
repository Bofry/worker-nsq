package internal

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/trace"
)

var _ MessageHandler = new(MessageHandleProc)

type MessageHandleProc func(ctx *Context, message *Message) error

func (proc MessageHandleProc) ProcessMessage(ctx *Context, message *Message) error {
	return proc(ctx, message)
}

var _ MessageHandleProc = StopRecursiveForwardMessageHandler

func StopRecursiveForwardMessageHandler(ctx *Context, msg *Message) error {
	ctx.logger.Fatal("invalid forward; it might be recursive forward message to InvalidMessageHandler")
	return nil
}

var (
	_ context.Context    = new(Context)
	_ trace.ValueContext = new(Context)
)

type Context struct {
	Channel string

	consumer *nsq.Consumer

	context        context.Context // parent context
	logger         *log.Logger
	disableLogging bool

	invalidMessageHandler MessageHandler
	invalidMessageSent    int32

	values     map[interface{}]interface{}
	valuesOnce sync.Once
}

// Deadline implements context.Context.
func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (*Context) Done() <-chan struct{} {
	return nil
}

// Err implements context.Context.
func (*Context) Err() error {
	return nil
}

// Value implements context.Context.
func (c *Context) Value(key interface{}) interface{} {
	if key == nil {
		return nil
	}
	if c.values != nil {
		v := c.values[key]
		if v != nil {
			return v
		}
	}
	if c.context != nil {
		return c.context.Value(key)
	}
	return nil
}

// SetValue implements trace.ValueContext.
func (c *Context) SetValue(key, value interface{}) {
	if key == nil {
		return
	}
	if c.values == nil {
		c.valuesOnce.Do(func() {
			if c.values == nil {
				c.values = make(map[interface{}]interface{})
			}
		})
	}
	c.values[key] = value
}

func (c *Context) Logger() *log.Logger {
	return c.logger
}

func (c *Context) IsRecordingLog() bool {
	return !c.disableLogging
}

func (c *Context) RecordingLog(v bool) {
	c.disableLogging = !v
}

func (c *Context) InvalidMessage(message *Message) error {
	if !atomic.CompareAndSwapInt32(&c.invalidMessageSent, 0, 1) {
		c.logger.Fatal("invalid operation; message has already been sent to InvalidMessageHandler")
	}

	GlobalContextHelper.InjectReplyCode(c, ABORT)

	if c.invalidMessageHandler != nil {
		var (
			tr       = GlobalTracerManager.GenerateManagedTracer(c.invalidMessageHandler)
			prevSpan = trace.SpanFromContext(c)
		)

		sp := tr.Start(prevSpan.Context(), __INVALID_MESSAGE_SPAN_NAME)
		defer sp.End()

		ctx := &Context{
			logger:                c.logger,
			values:                c.values,
			context:               c,
			invalidMessageHandler: MessageHandleProc(StopRecursiveForwardMessageHandler),
		}
		trace.SpanToContext(ctx, sp)

		err := c.invalidMessageHandler.ProcessMessage(ctx, message)
		_ = err // we won't process the error on InvalidMessageHandler
	}
	return nil
}

func (c *Context) Pause(topic ...string) error {
	return c.consumer.Pause(topic...)
}

func (c *Context) Resume(topic ...string) error {
	return c.consumer.Resume(topic...)
}

func (c *Context) Status() StatusCode {
	return GlobalContextHelper.ExtractReplyCode(c)
}

func (c *Context) clone() *Context {
	return &Context{
		Channel:               c.Channel,
		consumer:              c.consumer,
		context:               c.context,
		logger:                c.logger,
		invalidMessageHandler: c.invalidMessageHandler,
		values:                c.values,
	}
}
