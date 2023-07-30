package nsq

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-nsq/internal/middleware"
)

func UseErrorHandler(handler ErrorHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseLogging(service LoggingService) host.Middleware {
	if service == nil {
		panic("argument 'service' cannot be nil")
	}

	return &middleware.LoggingMiddleware{
		LoggingService: service,
	}
}

func UseMessageManager(messageManager interface{}) host.Middleware {
	if messageManager == nil {
		panic("argument 'messageManager' cannot be nil")
	}

	return &middleware.MessageManagerMiddleware{
		MessageManager: messageManager,
	}
}

func UseTracing(enabled bool) host.Middleware {
	return &middleware.TracingMiddleware{
		Enabled: enabled,
	}
}
