package internal

type NsqWorkerPreparer struct {
	worker *NsqWorker
}

func NewNsqWorkerPreparer(worker *NsqWorker) *NsqWorkerPreparer {
	return &NsqWorkerPreparer{
		worker: worker,
	}
}

func (p *NsqWorkerPreparer) RegisterUnhandledMessageHandler(handler MessageHandler) {
	p.worker.dispatcher.unhandledMessageHandler = handler
}

func (p *NsqWorkerPreparer) Router() Router {
	return p.worker.dispatcher.router
}
