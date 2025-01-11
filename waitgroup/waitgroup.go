//go:build !solution

package waitgroup

const initialCapacity = 1

type WaitGroup struct {
	runner  chan string
	counter chan int
}

func New() *WaitGroup {
	wg := &WaitGroup{
		runner:  nil,
		counter: make(chan int, initialCapacity),
	}
	wg.counter <- 0

	return wg
}

func (wg *WaitGroup) Add(delta int) {
	count := <-wg.counter
	defer wg.addDelta(count + delta)

	if count == 0 {
		wg.runner = make(chan string, initialCapacity)
	}
}

func (wg *WaitGroup) Done() {
	delta := <-wg.counter - 1
	defer wg.addDelta(delta)
	if delta == 0 {
		close(wg.runner)
	}

}

func (wg *WaitGroup) addDelta(delta int) {
	if delta < 0 {
		wg.counter <- 0
		wg.runner = nil
		panic("negative WaitGroup counter")
	}
	wg.counter <- delta
}

func (wg *WaitGroup) Wait() {
	if wg.runner == nil {
		return
	}
	<-wg.runner
}
