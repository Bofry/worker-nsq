package test

import (
	"log"

	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.LoggingService = new(LoggingService)

type LoggingService struct {
	logger *log.Logger
}

// ConfigureLogger implements middleware.LoggingService.
func (s *LoggingService) ConfigureLogger(l *log.Logger) {
	s.logger = l
}

// CreateEventLog implements middleware.LoggingService.
func (s *LoggingService) CreateEventLog(ev nsq.EventEvidence) nsq.EventLog {
	s.logger.Println("CreateEventLog()")
	return EventLog{
		logger:   s.logger,
		evidence: ev,
	}
}
