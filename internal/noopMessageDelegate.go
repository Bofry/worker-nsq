package internal

import (
	"time"

	nsq "github.com/Bofry/lib-nsq"
)

var _ nsq.MessageDelegate = NoopMessageDelegate(0)

type NoopMessageDelegate int

// OnFinish implements nsq.MessageDelegate.
func (NoopMessageDelegate) OnFinish(*nsq.Message) {}

// OnRequeue implements nsq.MessageDelegate.
func (NoopMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {}

// OnTouch implements nsq.MessageDelegate.
func (NoopMessageDelegate) OnTouch(*nsq.Message) {}
