package middleware

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
)

var _ EventLog = NoopEventLog(0)

type NoopEventLog int

// Flush implements EventLog.
func (n NoopEventLog) Flush() {}

// OnError implements EventLog.
func (n NoopEventLog) OnError(message *nsq.Message, err interface{}, stackTrace []byte) {}

// OnProcessMessage implements EventLog.
func (n NoopEventLog) OnProcessMessage(message *nsq.Message) {}

// OnProcessMessageComplete implements EventLog.
func (n NoopEventLog) OnProcessMessageComplete(message *nsq.Message, reply internal.ReplyCode) {}
