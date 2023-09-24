package internal

import (
	"reflect"

	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/trace"
)

var _ MessageObserver = CompositeMessageObserver(nil)

type CompositeMessageObserver []MessageObserver

// OnFinish implements MessageObserver.
func (o CompositeMessageObserver) OnFinish(ctx *Context, message *nsq.Message) {
	clonedMessage := message.Clone()
	clonedMessage.Delegate = GlobalRestrictedMessageDelegate

	var (
		sp        = trace.SpanFromContext(ctx)
		clonedCtx = ctx.clone()
	)
	clonedCtx.context = sp.Context()
	clonedCtx.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_InvalidOperation)

	for _, handler := range o {
		handler.OnFinish(clonedCtx, clonedMessage)
	}
}

// OnRequeue implements MessageObserver.
func (o CompositeMessageObserver) OnRequeue(ctx *Context, message *nsq.Message) {
	clonedMessage := message.Clone()
	clonedMessage.Delegate = GlobalRestrictedMessageDelegate

	var (
		sp        = trace.SpanFromContext(ctx)
		clonedCtx = ctx.clone()
	)
	clonedCtx.context = sp.Context()
	clonedCtx.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_InvalidOperation)

	for _, handler := range o {
		handler.OnRequeue(clonedCtx, clonedMessage)
	}
}

// OnTouch implements MessageObserver.
func (o CompositeMessageObserver) OnTouch(ctx *Context, message *nsq.Message) {
	clonedMessage := message.Clone()
	clonedMessage.Delegate = GlobalRestrictedMessageDelegate

	var (
		sp        = trace.SpanFromContext(ctx)
		clonedCtx = ctx.clone()
	)
	clonedCtx.context = sp.Context()
	clonedCtx.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_InvalidOperation)

	for _, handler := range o {
		handler.OnTouch(clonedCtx, clonedMessage)
	}
}

// Type implements MessageObserver.
func (o CompositeMessageObserver) Type() reflect.Type {
	return nil
}
