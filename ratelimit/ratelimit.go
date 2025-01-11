//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

const initialCapacity = 1

type Limiter struct {
	interval time.Duration
	timeouts chan []*time.Timer
	stop     chan string
	maxCount int
}

var ErrStopped = errors.New("limiter stopped")

func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	timeouts := getTimers(maxCount)
	timeoutsChan := getTimeoutsChan(timeouts)
	stopChan := getStopChan()

	return &Limiter{
		interval: interval,
		timeouts: timeoutsChan,
		stop:     stopChan,
		maxCount: maxCount,
	}
}

func getTimers(maxCount int) []*time.Timer {
	timeouts := make([]*time.Timer, maxCount)
	for i := range timeouts {
		timeouts[i] = time.NewTimer(0)
	}
	return timeouts
}

func getTimeoutsChan(timeouts []*time.Timer) chan []*time.Timer {
	timeoutsChan := make(chan []*time.Timer, initialCapacity)
	timeoutsChan <- timeouts
	return timeoutsChan
}

func getStopChan() chan string {
	return make(chan string, initialCapacity)
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stop:
		return ErrStopped
	default:
		return l.tryToUpdate(ctx)
	}
}

func (l *Limiter) tryToUpdate(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return l.update(ctx)
	}
}

func (l *Limiter) update(ctx context.Context) error {
	for i := 0; ; i = (i + 1) % l.maxCount {
		select {
		case <-l.stop:
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		case timeouts := <-l.timeouts:
			if l.isUpdated(timeouts, i) {
				return nil
			}
		}
	}
}

func (l *Limiter) isUpdated(timeouts []*time.Timer, i int) bool {
	defer func() {
		l.timeouts <- timeouts
	}()

	select {
	case <-timeouts[i].C:
		timeouts[i] = time.NewTimer(l.interval)
		return true

	default:
		return false
	}
}

func (l *Limiter) Stop() {
	close(l.stop)
}
