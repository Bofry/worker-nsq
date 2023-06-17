package middleware

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/reflecting"
	"github.com/Bofry/structproto/tagresolver"
	"github.com/Bofry/worker-nsq/internal"
)

var _ structproto.StructBinder = new(MessageManagerBinder)

type MessageManagerBinder struct {
	registrar *internal.NsqWorkerRegistrar
	app       *host.AppModule
}

func (b *MessageManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvMessageHandler := reflecting.AssignZero(rv)
	binder := &MessageHandlerBinder{
		messageHandlerType: rv.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.app.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.app.ServiceProvider(),
		},
	}
	err := b.bindMessageHandler(rvMessageHandler, binder)
	if err != nil {
		return err
	}

	// register MessageHandlers
	var (
		moduleID = field.IDName()
		topic    = field.Name()
	)
	return b.registerRoute(moduleID, topic, rvMessageHandler)
}

func (b *MessageManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageManagerBinder) bindMessageHandler(target reflect.Value, binder *MessageHandlerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *MessageManagerBinder) registerRoute(moduleID, topic string, rv reflect.Value) error {
	// register MessageHandlers
	if isMessageHandler(rv) {
		handler := asMessageHandler(rv)
		if handler != nil {
			if topic == UNHANDLED_MESSAGE_HANDLER_TOPIC_SYMBOL {
				// TODO: lack of moduleID ?
				b.registrar.SetUnhandledMessageHandler(handler)
			} else {
				// TODO: validate topic name comply NSQ topic spec
				b.registrar.AddRouter(topic, handler, moduleID)
			}
		}
	}
	return nil
}
