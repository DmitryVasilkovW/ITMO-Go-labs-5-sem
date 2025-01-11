//go:build !solution

package once

type Once struct {
	runner chan struct{}
	done   chan struct{}
}

func New() *Once {
	return &Once{
		runner: make(chan struct{}, 1),
		done:   make(chan struct{}, 1),
	}
}

func (o *Once) Do(f func()) {
	select {
	case o.runner <- struct{}{}:
		defer close(o.done)
		f()
	case <-o.done:
	}
}
