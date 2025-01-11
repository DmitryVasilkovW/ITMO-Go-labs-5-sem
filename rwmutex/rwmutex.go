//go:build !solution

package rwmutex

const initialCounter = 0
const InitialCapacity = 1

type RWMutex struct {
	reader chan int
	writer chan string
}

func New() *RWMutex {
	rw := &RWMutex{
		reader: make(chan int, InitialCapacity),
		writer: make(chan string, InitialCapacity),
	}
	rw.reader <- initialCounter
	rw.writer <- "239"

	return rw
}

func (rw *RWMutex) RLock() {
	counter := <-rw.reader
	if counter == 0 {
		<-rw.writer
	}

	rw.reader <- counter + 1
}

func (rw *RWMutex) RUnlock() {
	counter := <-rw.reader
	if counter == 1 {
		rw.writer <- "239"
	}

	rw.reader <- counter - 1
}

func (rw *RWMutex) Lock() {
	<-rw.writer
}

func (rw *RWMutex) Unlock() {
	rw.writer <- "239"
}
