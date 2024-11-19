package middleware

import (
	"context"
	"runtime/debug"

	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
)

var _ internal.MessageHandleModule = new(LoggingHandleModule)

type LoggingHandleModule struct {
	successor      internal.MessageHandleModule
	loggingService LoggingService
}

// CanSetSuccessor implements internal.MessageHandleModule.
func (*LoggingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements internal.MessageHandleModule.
func (m *LoggingHandleModule) SetSuccessor(successor internal.MessageHandleModule) {
	m.successor = successor
}

// ProcessMessage implements internal.MessageHandleModule.
func (m *LoggingHandleModule) ProcessMessage(ctx *internal.Context, message *nsq.Message, state internal.ProcessingState, recover *internal.Recover) error {
	if m.successor != nil {
		if !ctx.IsRecordingLog() {
			return nil
		}

		evidence := EventEvidence{
			traceID: state.Span.TraceID(),
			spanID:  state.Span.SpanID(),
			topic:   state.Topic,
		}

		eventLog := m.loggingService.CreateEventLog(evidence)
		// NOTE restrict call Finish(), Requeue(), Touch()
		internal.GlobalMessageDelegateHelper.Restrict(message)
		eventLog.OnProcessMessage(message)

		return recover.
			Defer(func(err interface{}) {
				if err != nil {
					defer func() {
						// NOTE restrict call Finish(), Requeue(), Touch()
						internal.GlobalMessageDelegateHelper.Restrict(message)
						eventLog.OnError(message, err, debug.Stack())
						eventLog.Flush()
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					var (
						reply internal.ReplyCode = internal.GlobalContextHelper.ExtractReplyCode(ctx)
					)
					defer eventLog.Flush()

					// NOTE restrict call Finish(), Requeue(), Touch()
					internal.GlobalMessageDelegateHelper.Restrict(message)
					eventLog.OnProcessMessageComplete(message, reply)
				}
			}).
			Do(func(f internal.Finalizer) error {
				// NOTE restrict call Finish(), Requeue(), Touch()
				internal.GlobalMessageDelegateHelper.Restrict(message)
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
