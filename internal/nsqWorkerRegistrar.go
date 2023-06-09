package internal

type NsqWorkerRegistrar struct {
	worker *NsqWorker
}

func NewNsqWorkerRegistrar(worker *NsqWorker) *NsqWorkerRegistrar {
	return &NsqWorkerRegistrar{
		worker: worker,
	}
}

func (r *NsqWorkerRegistrar) RegisterMessageHandleModule(module MessageHandleModule) {
	r.worker.messageHandleService.Register(module)
}

func (r *NsqWorkerRegistrar) EnableTracer(enabled bool) {
	r.worker.messageTracerService.Enabled = enabled
}

func (r *NsqWorkerRegistrar) SetErrorHandler(handler ErrorHandler) {
	r.worker.messageDispatcher.ErrorHandler = handler
}

func (r *NsqWorkerRegistrar) SetUnhandledMessageHandler(handler MessageHandler) {
	r.worker.messageDispatcher.UnhandledMessageHandler = handler
}

func (r *NsqWorkerRegistrar) SetMessageManager(messageManager interface{}) {
	r.worker.messageManager = messageManager
}

func (r *NsqWorkerRegistrar) AddRouter(topic string, handler MessageHandler, handlerComponentID string) {
	r.worker.router.Add(topic, handler, handlerComponentID)
}
