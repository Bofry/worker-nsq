package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/worker-nsq/internal"
)

var _ host.Middleware = new(LoggingMiddleware)

type LoggingMiddleware struct {
	LoggingService LoggingService
}

// Init implements internal.Middleware
func (m *LoggingMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asNsqWorker(app.Host())
		registrar = NewNsqWorkerRegistrar(worker)
	)

	m.LoggingService.ConfigureLogger(worker.Logger())

	loggingHandleModule := &LoggingHandleModule{
		loggingService: m.LoggingService,
	}
	registrar.RegisterMessageHandleModule(loggingHandleModule)
}
