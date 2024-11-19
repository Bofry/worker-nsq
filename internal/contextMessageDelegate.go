package internal

import (
	"sync"
	"sync/atomic"
	"time"

	nsq "github.com/Bofry/lib-nsq"
)

var _ nsq.MessageDelegate = new(ContextMessageDelegate)

type ContextMessageDelegate struct {
	parent nsq.MessageDelegate

	ctx *Context

	messageObserver MessageObserver

	restricted int32
	mu         sync.Mutex
}

func NewContextMessageDelegate(ctx *Context) *ContextMessageDelegate {
	return &ContextMessageDelegate{
		ctx: ctx,
	}
}

func (d *ContextMessageDelegate) OnFinish(msg *nsq.Message) {
	if d.isRestricted() {
		GlobalNoopMessageDelegate.OnFinish(nil)
		return
	}

	d.parent.OnFinish(msg)
	GlobalContextHelper.InjectReplyCodeSafe(d.ctx, PASS)

	// observer
	if d.messageObserver != nil {
		d.messageObserver.OnFinish(d.ctx, msg)
	}
}

func (d *ContextMessageDelegate) OnRequeue(msg *nsq.Message, delay time.Duration, backoff bool) {
	if d.isRestricted() {
		GlobalNoopMessageDelegate.OnRequeue(nil, delay, backoff)
		return
	}

	d.parent.OnRequeue(msg, delay, backoff)
	GlobalContextHelper.InjectReplyCodeSafe(d.ctx, FAIL)

	// observer
	if d.messageObserver != nil {
		d.messageObserver.OnRequeue(d.ctx, msg)
	}
}

func (d *ContextMessageDelegate) OnTouch(msg *nsq.Message) {
	if d.isRestricted() {
		GlobalNoopMessageDelegate.OnTouch(nil)
		return
	}

	d.parent.OnTouch(msg)

	// observer
	if d.messageObserver != nil {
		d.messageObserver.OnTouch(d.ctx, msg)
	}
}

func (d *ContextMessageDelegate) configure(msg *nsq.Message) {
	if d.parent == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if d.parent == nil {
			d.parent = msg.Delegate
			msg.Delegate = d
		}
	}
}

func (d *ContextMessageDelegate) isRestricted() bool {
	return atomic.LoadInt32(&d.restricted) == 1
}

func (d *ContextMessageDelegate) restrict() {
	atomic.StoreInt32(&d.restricted, 1)
}

func (d *ContextMessageDelegate) unrestrict() {
	atomic.StoreInt32(&d.restricted, 0)
}

func (d *ContextMessageDelegate) registerMessageObservers(observers []MessageObserver) {
	d.messageObserver = CompositeMessageObserver(observers)
}

func (d *ContextMessageDelegate) unregisterAllMessageObservers() {
	d.messageObserver = nil
}
