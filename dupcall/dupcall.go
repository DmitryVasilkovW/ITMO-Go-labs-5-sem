//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	mu        sync.Mutex
	ch        chan struct{}
	resources interface{}
	err       error
}

func (o *Call) Do(ctx context.Context, cb func(context.Context) (interface{}, error)) (interface{}, error) {
	o.runOnce(ctx, cb)
	return o.waitForCompletion(ctx)
}

func (o *Call) waitForCompletion(ctx context.Context) (interface{}, error) {
	o.mu.Lock()
	nch := o.ch
	o.mu.Unlock()

	select {
	case <-nch:
		o.mu.Lock()
		defer o.mu.Unlock()
		return o.resources, o.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (o *Call) runOnce(ctx context.Context, cb func(context.Context) (interface{}, error)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.ch == nil {
		o.ch = make(chan struct{})
		go o.start(ctx, cb)
	}
}

func (o *Call) start(ctx context.Context, cb func(context.Context) (interface{}, error)) {
	o.setResourcesAndError(ctx, cb)
	o.closeChannel()
}

func (o *Call) setResourcesAndError(ctx context.Context, cb func(context.Context) (interface{}, error)) {
	o.resources, o.err = cb(ctx)
}

func (o *Call) closeChannel() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.ch != nil {
		close(o.ch)
		o.ch = nil
	}
}
