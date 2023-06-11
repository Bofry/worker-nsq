package internal

import (
	"fmt"
	"io"
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostModule = NsqWorkerModule{}

type NsqWorkerModule struct{}

// ConfigureLogger implements host.HostService
func (NsqWorkerModule) ConfigureLogger(logflags int, w io.Writer) {
	fmt.Println("******NsqWorkerModule.ConfigureLogger()")
	NsqWorkerLogger.SetFlags(logflags)
	NsqWorkerLogger.SetOutput(w)
}

// Init implements host.HostService
func (NsqWorkerModule) Init(h host.Host, app *host.AppModule) {
	if v, ok := h.(*NsqWorker); ok {
		v.alloc()
		v.setTracerProvider(app.TracerProvider())
		v.setTextMapPropagator(app.TextMapPropagator())
		v.setLogger(app.Logger())
	}
}

// InitComplete implements host.HostService
func (NsqWorkerModule) InitComplete(h host.Host, app *host.AppModule) {
	if v, ok := h.(*NsqWorker); ok {
		v.init()
	}
}

// DescribeHostType implements host.HostService
func (NsqWorkerModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
