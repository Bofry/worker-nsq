package middleware

import (
	"context"
	"runtime/debug"

	nsq "github.com/Bofry/lib-nsq"
	. "github.com/Bofry/worker-nsq/internal"
)

var _ MessageHandleModule = new(LoggingHandleModule)

type LoggingHandleModule struct {
	successor      MessageHandleModule
	loggingService LoggingService
}

// CanSetSuccessor implements internal.MessageHandleModule.
func (*LoggingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements internal.MessageHandleModule.
func (m *LoggingHandleModule) SetSuccessor(successor MessageHandleModule) {
	m.successor = successor
}

// ProcessMessage implements internal.MessageHandleModule.
func (m *LoggingHandleModule) ProcessMessage(ctx *Context, message *nsq.Message, state ProcessingState, recover *Recover) error {
	if m.successor != nil {
		evidence := EventEvidence{
			traceID: state.Span.TraceID(),
			spanID:  state.Span.SpanID(),
			topic:   state.Topic,
		}

		eventLog := m.loggingService.CreateEventLog(evidence)
		eventLog.BeforeProcessMessage(message)

		return recover.
			Defer(func(err interface{}) {
				if err != nil {
					defer func() {
						eventLog.LogError(message, err, debug.Stack())
						eventLog.Flush()
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					eventLog.AfterProcessMessage(message)
					eventLog.Flush()
				}
			}).
			Do(func(f Finalizer) error {
				return m.successor.ProcessMessage(ctx, message, state, recover)
			})
	}
	return nil
}

// OnInitComplete implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	return nil
}

// OnStop implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}
