package middleware

import "log"

var _ LoggingService = NoopLoggingService{}

type NoopLoggingService struct{}

// ConfigureLogger implements LoggingService.
func (n NoopLoggingService) ConfigureLogger(l *log.Logger) {}

// CreateEventLog implements LoggingService.
func (n NoopLoggingService) CreateEventLog(ev EventEvidence) EventLog {
	return NoopEventLogSingleton
}
