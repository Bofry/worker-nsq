package internal

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Bofry/trace"
)

var _ MessageHandler = new(MessageHandleProc)

type MessageHandleProc func(ctx *Context, message *Message) error

func (proc MessageHandleProc) ProcessMessage(ctx *Context, message *Message) error {
	return proc(ctx, message)
}

var _ MessageHandleProc = StopRecursiveForwardMessageHandler

func StopRecursiveForwardMessageHandler(ctx *Context, msg *Message) error {
	ctx.logger.Fatal("invalid forward; it might be recursive forward message to UnhandledMessageHandler")
	return nil
}

var (
	_ context.Context    = new(Context)
	_ trace.ValueContext = new(Context)
)

type Context struct {
	Channel string

	logger *log.Logger

	unhandledMessageHandler MessageHandler

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
		return c.values[key]
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

func (c *Context) ForwardUnhandledMessage(message *Message) error {
	if c.unhandledMessageHandler != nil {
		ctx := &Context{
			logger:                  c.logger,
			unhandledMessageHandler: MessageHandleProc(StopRecursiveForwardMessageHandler),
		}
		return c.unhandledMessageHandler.ProcessMessage(ctx, message)
	}
	return nil
}
