package internal

import (
	"context"
	"fmt"

	"github.com/Bofry/lib-nsq/tracing"
	"github.com/Bofry/trace"
)

var _ MessageHandler = new(MessageDispatcher)

type MessageDispatcher struct {
	MessageHandleService *MessageHandleService
	MessageTracerService *MessageTracerService
	Router               Router

	OnHostErrorProc OnHostErrorHandler

	ErrorHandler          ErrorHandler
	InvalidMessageHandler MessageHandler
}

func (d *MessageDispatcher) Topics() []string {
	var (
		router = d.Router
	)

	if router != nil {
		keys := make([]string, 0, len(router))
		for k := range router {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *MessageDispatcher) ProcessMessage(ctx *Context, message *Message) error {
	// start tracing
	var (
		componentID = d.Router.FindHandlerComponentID(message.Topic)
		carrier     = tracing.NewMessageStateCarrier(&message.Content().State)

		spanName string = message.Topic
		tr       *trace.SeverityTracer
		sp       *trace.SeveritySpan
	)

	tr = d.MessageTracerService.Tracer(componentID)
	sp = tr.ExtractWithPropagator(
		ctx,
		d.MessageTracerService.TextMapPropagator(),
		carrier,
		spanName)
	defer sp.End()

	processingState := ProcessingState{
		Topic:  message.Topic,
		Tracer: tr,
		Span:   sp,
	}

	// set invalidMessageHandler
	ctx.invalidMessageHandler = d.InvalidMessageHandler

	return d.MessageHandleService.ProcessMessage(ctx, message, processingState, new(Recover))
}

func (d *MessageDispatcher) internalProcessMessage(ctx *Context, message *Message, state ProcessingState, recover *Recover) error {
	return recover.
		Defer(func(err interface{}) {
			if err != nil {
				d.processError(ctx, message, err)
			}
		}).
		Do(func(finalizer Finalizer) error {
			var (
				tr    *trace.SeverityTracer = state.Tracer
				sp    *trace.SeveritySpan   = state.Span
				topic string                = state.Topic
			)
			_ = tr

			// set Span
			trace.SpanToContext(ctx, sp)

			finalizer.Add(func(err interface{}) {
				if err != nil {
					if e, ok := err.(error); ok {
						sp.Err(e)
					} else if e, ok := err.(string); ok {
						sp.Err(fmt.Errorf(e))
					} else if e, ok := err.(fmt.Stringer); ok {
						sp.Err(fmt.Errorf(e.String()))
					} else {
						sp.Err(fmt.Errorf("%+v", err))
					}
				}

				var (
					reply = GlobalContextHelper.ExtractReplyCode(ctx)
				)

				switch reply {
				case PASS:
					sp.Reply(trace.PASS, reply)
				case FAIL, ABORT:
					sp.Reply(trace.FAIL, reply)
				}
			})

			sp.Tags(
				// TODO: add nsq server version
				trace.Topic(topic),
				trace.ConsumerGroup(ctx.Channel),
				trace.BrokerIP(message.NSQDAddress),
				trace.MessageID(string(message.ID[:])),
				trace.Key(__ATTR_ATTEMPTS).Int(int(message.Attempts)),
			)

			handler := d.Router.Get(topic)
			if handler != nil {
				err := handler.ProcessMessage(ctx, message)
				if err == nil {
					reply := GlobalContextHelper.ExtractReplyCode(ctx)
					if reply == UNSET {
						message.Finish()
					}
				}
				return err
			}
			return ctx.ThrowInvalidMessageError(message)
		})
}

func (d *MessageDispatcher) init() {
	// register the default MessageHandleModule
	stdMessageHandleModule := NewStdMessageHandleModule(d)
	d.MessageHandleService.Register(stdMessageHandleModule)
}

func (d *MessageDispatcher) processError(ctx *Context, message *Message, err interface{}) {
	if d.ErrorHandler != nil {
		d.ErrorHandler(ctx, message, err)
	}
}

func (d *MessageDispatcher) start(ctx context.Context) {
	err := d.MessageHandleService.triggerStart(ctx)
	if err != nil {
		var disposed bool = false
		if d.OnHostErrorProc != nil {
			disposed = d.OnHostErrorProc(err)
		}
		if !disposed {
			NsqWorkerLogger.Fatalf("%+v", err)
		}
	}
}

func (d *MessageDispatcher) stop(ctx context.Context) {
	for err := range d.MessageHandleService.triggerStop(ctx) {
		if err != nil {
			NsqWorkerLogger.Printf("%+v", err)
		}
	}
}
