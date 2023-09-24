package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/tagresolver"

	. "github.com/Bofry/worker-nsq/internal"
)

var _ host.Middleware = new(MessageObserverManagerMiddleware)

type MessageObserverManagerMiddleware struct {
	MessageObserverManager interface{}
}

// Init implements internal.Middleware.
func (m *MessageObserverManagerMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asNsqWorker(app.Host())
		registrar = NewNsqWorkerRegistrar(worker)
	)

	// binding MessageObserverManager
	binder := &MessageObserverManagerBinder{
		registrar: registrar,
		app:       app,
	}

	err := m.bindMessageObserverManager(m.MessageObserverManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *MessageObserverManagerMiddleware) bindMessageObserverManager(target interface{}, binder *MessageObserverManagerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
