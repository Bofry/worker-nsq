package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/worker-nsq/internal"
)

var _ host.Middleware = new(TracingMiddleware)

type TracingMiddleware struct {
	Enabled bool
}

// Init implements internal.Middleware.
func (m *TracingMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asNsqWorker(app.Host())
		registrar = NewNsqWorkerRegistrar(worker)
	)

	registrar.EnableTracer(m.Enabled)
}
