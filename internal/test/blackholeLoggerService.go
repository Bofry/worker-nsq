package test

import (
	"bytes"
	"log"

	nsq "github.com/Bofry/worker-nsq"
)

var _ nsq.LoggingService = new(BlackholeLoggerService)

type BlackholeLoggerService struct {
	Buffer *bytes.Buffer
}

func (s *BlackholeLoggerService) CreateEventLog(ev nsq.EventEvidence) nsq.EventLog {
	s.Buffer.WriteString("CreateEventLog()")
	s.Buffer.WriteByte('\n')
	return &BlackholeEventLog{
		buffer: s.Buffer,
	}
}

func (*BlackholeLoggerService) ConfigureLogger(l *log.Logger) {
}

var _ nsq.EventLog = new(BlackholeEventLog)

type BlackholeEventLog struct {
	buffer *bytes.Buffer
}

func (l *BlackholeEventLog) LogError(message *nsq.Message, err interface{}, stackTrace []byte) {
	l.buffer.WriteString("LogError()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) AfterProcessMessage(message *nsq.Message) {
	l.buffer.WriteString("AfterProcessMessage()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) BeforeProcessMessage(message *nsq.Message) {
	l.buffer.WriteString("BeforeProcessMessage()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) Flush() {
	l.buffer.WriteString("Flush()")
	l.buffer.WriteByte('\n')
}
