//go:build !solution

package cond

const capacity = 1

type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L    Locker
	chin chan chan string
}

func New(l Locker) *Cond {
	return &Cond{L: l, chin: make(chan chan string, 1000)}
}

func (c *Cond) Wait() {
	defer func() {
		c.L.Lock()
	}()

	ch := make(chan string, capacity)

	c.chin <- ch
	c.L.Unlock()
	<-ch
}

func (c *Cond) Signal() {
	c.sendSignal()
}

func (c *Cond) Broadcast() {
	flag := c.sendSignal()
	for flag {
		flag = c.sendSignal()
	}
}

func (c *Cond) sendSignal() bool {
	select {
	case ch := <-c.chin:
		return hasNotified(ch)
	default:
	}
	return false
}

func hasNotified(ch chan string) bool {
	select {
	case ch <- "🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓":
	default:

	}

	return true
}
