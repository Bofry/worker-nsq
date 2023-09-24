package test

import (
	"fmt"
	"reflect"

	nsq "github.com/Bofry/worker-nsq"
	"github.com/Bofry/worker-nsq/tracing"
)

var _ nsq.MessageObserver = new(GoTestStreamMessageObserver)

type GoTestStreamMessageObserver struct {
	ServiceProvider *ServiceProvider
}

func (*GoTestStreamMessageObserver) Init() {
	fmt.Println("GoTestStreamMessageObserver.Init()")
}

// OnFinish implements internal.MessageObserver.
func (o *GoTestStreamMessageObserver) OnFinish(ctx *nsq.Context, message *nsq.Message) {
	tr := tracing.GetTracer(o)
	sp := tr.Start(ctx, "OnFinish()")
	defer sp.End()

	o.ServiceProvider.Logger().Println("GoTestStreamMessageObserver.OnFinish()")
}

// OnRequeue implements internal.MessageObserver.
func (o *GoTestStreamMessageObserver) OnRequeue(ctx *nsq.Context, message *nsq.Message) {
	tr := tracing.GetTracer(o)
	sp := tr.Start(ctx, "OnRequeue()")
	defer sp.End()

	o.ServiceProvider.Logger().Println("GoTestStreamMessageObserver.OnRequeue()")
}

// OnTouch implements internal.MessageObserver.
func (o *GoTestStreamMessageObserver) OnTouch(ctx *nsq.Context, message *nsq.Message) {
	o.ServiceProvider.Logger().Println("GoTestStreamMessageObserver.OnTouch()")
}

func (o *GoTestStreamMessageObserver) Type() reflect.Type {
	return reflect.TypeOf(o)
}
