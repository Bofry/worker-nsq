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

func (l *BlackholeEventLog) OnError(message *nsq.Message, err interface{}, stackTrace []byte) {
	l.buffer.WriteString("LogError()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessMessageComplete(message *nsq.Message, reply nsq.ReplyCode) {
	l.buffer.WriteString("OnProcessMessageComplete()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessMessage(message *nsq.Message) {
	l.buffer.WriteString("OnProcessMessage()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) Flush() {
	l.buffer.WriteString("Flush()")
	l.buffer.WriteByte('\n')
}
