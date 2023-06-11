package test

import (
	"log"

	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.LoggingService = new(MultiLoggerService)

type MultiLoggerService struct {
	LoggingServices []nsq.LoggingService
}

func (s *MultiLoggerService) CreateEventLog(ev nsq.EventEvidence) nsq.EventLog {
	var eventlogs []nsq.EventLog
	for _, svc := range s.LoggingServices {
		eventlogs = append(eventlogs, svc.CreateEventLog(ev))
	}

	return MultiEventLog{
		EventLogs: eventlogs,
	}
}

func (s *MultiLoggerService) ConfigureLogger(l *log.Logger) {
	for _, svc := range s.LoggingServices {
		svc.ConfigureLogger(l)
	}
}

var _ nsq.EventLog = MultiEventLog{}

type MultiEventLog struct {
	EventLogs []nsq.EventLog
}

// Flush implements middleware.EventLog.
func (l MultiEventLog) Flush() {
	for _, log := range l.EventLogs {
		log.Flush()
	}
}

// LogError implements middleware.EventLog.
func (l MultiEventLog) LogError(message *nsq.Message, err interface{}, stackTrace []byte) {
	for _, log := range l.EventLogs {
		log.LogError(message, err, stackTrace)
	}
}

// AfterProcessMessage implements middleware.EventLog.
func (l MultiEventLog) AfterProcessMessage(message *nsq.Message) {
	for _, log := range l.EventLogs {
		log.AfterProcessMessage(message)
	}
}

// BeforeProcessMessage implements middleware.EventLog.
func (l MultiEventLog) BeforeProcessMessage(message *nsq.Message) {
	for _, log := range l.EventLogs {
		log.BeforeProcessMessage(message)
	}
}
