package nsq

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-nsq/internal/middleware"
)

func UseTopicGateway(topicGateway interface{}) host.Middleware {
	if topicGateway == nil {
		panic("argument 'topicGateway' cannot be nil")
	}

	return &middleware.TopicGatewayMiddleware{
		TopicGateway: topicGateway,
	}
}
