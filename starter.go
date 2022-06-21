package nsq

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-nsq/internal"
)

func Startup(app interface{}) *host.Starter {
	var (
		starter = host.Startup(app)
	)

	host.RegisterHostService(starter, internal.NsqWorkerServiceInstance)

	return starter
}
