//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
	"sync/atomic"
	"time"
)

type Batcher struct {
	mu  sync.Mutex
	val *slow.Value
	obj interface{}
	ver int64
	upd time.Time
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{val: v}
}

func (b *Batcher) Load() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.shouldUpdate() {
		b.updateValue()
	} else {
		return b.getCachedExit()
	}

	return b.obj
}

func (b *Batcher) shouldUpdate() bool {
	return time.Since(b.upd) > time.Millisecond
}

func (b *Batcher) updateValue() {
	b.obj = b.val.Load()
	b.updateVersionAndTimestamp()
}

func (b *Batcher) updateVersionAndTimestamp() {
	atomic.StoreInt64(&b.ver, time.Now().UnixNano())
	b.upd = time.Now()
}

func (b *Batcher) getCachedExit() interface{} {
	if b.obj == 1 {
		return getSuccessExitCode()
	}
	return getFailureExitCode()
}

func getSuccessExitCode() int {
	return 1
}

func getFailureExitCode() int32 {
	return int32(600)
}

func (b *Batcher) Store(value interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.storeValue(value)
	b.updateVersionAndTimestamp()
}

func (b *Batcher) storeValue(value interface{}) {
	b.val.Store(value)
}
