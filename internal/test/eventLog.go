package test

import (
	"fmt"
	"log"

	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.EventLog = EventLog{}

type EventLog struct {
	logger   *log.Logger
	evidence nsq.EventEvidence
}

// AfterProcessMessage implements middleware.EventLog.
func (l EventLog) AfterProcessMessage(message *nsq.Message) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.ProcessingSpanID(),
		l.evidence.ProcessingSpanID())

	l.logger.Printf("EventLog.AfterProcessMessage(): (%s) %s\n", traceID, string(message.ID[:]))
}

// BeforeProcessMessage implements middleware.EventLog.
func (l EventLog) BeforeProcessMessage(message *nsq.Message) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.ProcessingSpanID(),
		l.evidence.ProcessingSpanID())

	l.logger.Printf("EventLog.BeforeProcessMessage(): (%s) %s\n", traceID, string(message.ID[:]))
}

// Flush implements middleware.EventLog.
func (l EventLog) Flush() {
	l.logger.Println("EventLog.Flush()")
}

// LogError implements middleware.EventLog.
func (l EventLog) LogError(message *nsq.Message, err interface{}, stackTrace []byte) {
	l.logger.Printf("EventLog.WriteError(): %v\n", err)
}
