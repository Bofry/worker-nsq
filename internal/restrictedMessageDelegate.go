package internal

import (
	"time"

	nsq "github.com/Bofry/lib-nsq"
)

var _ nsq.MessageDelegate = RestrictedMessageDelegate(0)

type RestrictedMessageDelegate int

// OnFinish implements nsq.MessageDelegate.
func (RestrictedMessageDelegate) OnFinish(*nsq.Message) {
	panic(RestrictedOperationError("Message.Finish() cannot be called by restricted area"))
}

// OnRequeue implements nsq.MessageDelegate.
func (RestrictedMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {
	panic(RestrictedOperationError("Message.Requeue() cannot be called by restricted area"))
}

// OnTouch implements nsq.MessageDelegate.
func (RestrictedMessageDelegate) OnTouch(*nsq.Message) {
	panic(RestrictedOperationError("Message.Touch() cannot be called by restricted area"))
}
