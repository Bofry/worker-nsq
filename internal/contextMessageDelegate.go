package internal

import (
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

var _ nsq.MessageDelegate = new(ContextMessageDelegate)

type ContextMessageDelegate struct {
	parent nsq.MessageDelegate

	ctx *Context

	mu sync.Mutex
}

func NewContextMessageDelegate(ctx *Context) *ContextMessageDelegate {
	return &ContextMessageDelegate{
		ctx: ctx,
	}
}

func (d *ContextMessageDelegate) OnFinish(msg *nsq.Message) {
	d.parent.OnFinish(msg)
	GlobalContextHelper.InjectReplyCode(d.ctx, PASS)
}

func (d *ContextMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {
	d.parent.OnRequeue(m, delay, backoff)
	GlobalContextHelper.InjectReplyCode(d.ctx, FAIL)
}

func (d *ContextMessageDelegate) OnTouch(msg *nsq.Message) {
	d.parent.OnTouch(msg)
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
