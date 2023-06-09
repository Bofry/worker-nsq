package internal

import (
	"context"

	nsq "github.com/Bofry/lib-nsq"
)

var _ MessageHandleModule = new(StdMessageHandleModule)

type StdMessageHandleModule struct {
	dispatcher *MessageDispatcher
}

func NewStdMessageHandleModule(dispatcher *MessageDispatcher) *StdMessageHandleModule {
	return &StdMessageHandleModule{
		dispatcher: dispatcher,
	}
}

// CanSetSuccessor implements MessageHandleModule.
func (*StdMessageHandleModule) CanSetSuccessor() bool {
	return false
}

// SetSuccessor implements MessageHandleModule.
func (*StdMessageHandleModule) SetSuccessor(successor MessageHandleModule) {
	panic("unsupported operation")
}

// ProcessMessage implements MessageHandleModule.
func (m *StdMessageHandleModule) ProcessMessage(ctx *Context, message *nsq.Message, state ProcessingState, recover *Recover) error {
	return m.dispatcher.internalProcessMessage(ctx, message, state, recover)
}

// OnInitComplete implements MessageHandleModule.
func (*StdMessageHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements MessageHandleModule.
func (*StdMessageHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	return nil
}

// OnStop implements MessageHandleModule.
func (*StdMessageHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}
