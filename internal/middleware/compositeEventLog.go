package middleware

import (
	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/worker-nsq/internal"
)

var _ EventLog = CompositeEventLog{}

type CompositeEventLog struct {
	eventLogs []EventLog
}

// Flush implements EventLog.
func (l CompositeEventLog) Flush() {
	for _, log := range l.eventLogs {
		log.Flush()
	}
}

// OnError implements EventLog.
func (l CompositeEventLog) OnError(message *nsq.Message, err interface{}, stackTrace []byte) {
	for _, log := range l.eventLogs {
		log.OnError(message, err, stackTrace)
	}
}

// OnProcessMessage implements EventLog.
func (l CompositeEventLog) OnProcessMessage(message *nsq.Message) {
	for _, log := range l.eventLogs {
		log.OnProcessMessage(message)
	}
}

// OnProcessMessageComplete implements EventLog.
func (l CompositeEventLog) OnProcessMessageComplete(message *nsq.Message, reply internal.ReplyCode) {
	for _, log := range l.eventLogs {
		log.OnProcessMessageComplete(message, reply)
	}
}
