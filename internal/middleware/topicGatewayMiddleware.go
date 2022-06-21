package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/worker-nsq/internal"
)

var _ host.Middleware = new(TopicGatewayMiddleware)

type TopicGatewayMiddleware struct {
	TopicGateway interface{}
}

func (m *TopicGatewayMiddleware) Init(appCtx *host.AppContext) {
	var (
		nsqworker = asNsqWorker(appCtx.Host())
		preparer  = internal.NewNsqWorkerPreparer(nsqworker)
	)

	binder := &TopicGatewayBinder{
		router:                           preparer.Router(),
		appContext:                       appCtx,
		configureUnhandledMessageHandler: preparer.RegisterUnhandledMessageHandler,
	}

	err := m.performBindTopicGateway(m.TopicGateway, binder)
	if err != nil {
		panic(err)
	}
}

func (m *TopicGatewayMiddleware) performBindTopicGateway(target interface{}, binder *TopicGatewayBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagName:     TAG_TOPIC,
			TagResolver: TopicTagResolve,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
