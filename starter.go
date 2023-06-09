package nsq

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-nsq/internal"
)

func Startup(app interface{}) *host.Starter {
	var (
		starter = host.Startup(app)
	)

	host.RegisterHostModule(starter, internal.NsqWorkerModuleInstance)

	return starter
}
