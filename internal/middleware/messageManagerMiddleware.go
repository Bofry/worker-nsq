package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	. "github.com/Bofry/worker-nsq/internal"
)

var _ host.Middleware = new(MessageManagerMiddleware)

type MessageManagerMiddleware struct {
	MessageManager interface{}
}

func (m *MessageManagerMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asNsqWorker(app.Host())
		registrar = NewNsqWorkerRegistrar(worker)
	)

	// register RequestManager offer FasthttpHost processing later.
	registrar.SetMessageManager(m.MessageManager)

	// binding MessageManage
	binder := &MessageManagerBinder{
		registrar: registrar,
		app:       app,
	}

	err := m.bindMessageManager(m.MessageManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *MessageManagerMiddleware) bindMessageManager(target interface{}, binder *MessageManagerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagName:             TAG_TOPIC,
			TagResolver:         TopicTagResolver,
			CheckDuplicateNames: true,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
