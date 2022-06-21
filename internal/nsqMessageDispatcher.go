package internal

type NsqMessageDispatcher struct {
	messageHandler          MessageHandleProc
	unhandledMessageHandler MessageHandler
	router                  Router
}

func NewNsqMessageDispatcher() *NsqMessageDispatcher {
	return &NsqMessageDispatcher{
		router: make(Router),
	}
}

func (d *NsqMessageDispatcher) Topics() []string {
	var (
		router = d.router
	)

	if router != nil {
		keys := make([]string, 0, len(router))
		for k := range router {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *NsqMessageDispatcher) ProcessMessage(message *Message) error {
	handler := d.router.Get(message.Topic)
	if handler != nil {
		return handler.ProcessMessage(message)
	}
	return message.ForwardUnhandledMessageHandler()
}

func (d *NsqMessageDispatcher) ProcessUnhandledMessage(message *Message) error {
	if d.unhandledMessageHandler != nil {
		// prevent recursive call ForwardUnhandledMessageHandler
		message.StopForwardUnhandledMessageHandler()
		return d.unhandledMessageHandler.ProcessMessage(message)
	}
	return nil
}
