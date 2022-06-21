package internal

import (
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostService = new(NsqWorkerService)

type NsqWorkerService struct{}

func (p *NsqWorkerService) Init(h host.Host, ctx *host.AppContext) {
	if v, ok := h.(*NsqWorker); ok {
		v.preInit()
	}
}

func (p *NsqWorkerService) InitComplete(h host.Host, ctx *host.AppContext) {
	if v, ok := h.(*NsqWorker); ok {
		v.init()
	}
}

func (p *NsqWorkerService) GetHostType() reflect.Type {
	return typeOfHost
}
