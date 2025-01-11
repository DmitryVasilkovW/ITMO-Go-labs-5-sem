package externalsort

import (
	"container/heap"
	"errors"
	"io"
)

func initHeap(readers []LineReader) (Heap, error) {
	h := make(Heap, 0, len(readers))

	for _, reader := range readers {
		line, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) || (errors.Is(err, io.EOF) && len(line) == 0) {
			continue
		}

		heap.Push(&h, &HeapItem{
			lr:  &reader,
			top: line,
		})
	}

	return h, nil
}

func mergeLines(w LineWriter, h Heap) error {
	for h.Len() > 0 {
		item := heap.Pop(&h).(*HeapItem)

		err := w.Write(item.top)
		if err != nil {
			return err
		}

		err = pushToHeap(item, &h)
		if err != nil {
			return err
		}
	}
	return nil
}

func pushToHeap(item *HeapItem, h *Heap) error {
	line, err := (*item.lr).ReadLine()
	if err == nil || (errors.Is(err, io.EOF) && len(line) > 0) {
		heap.Push(h, &HeapItem{
			lr:  item.lr,
			top: line,
		})
	}
	return nil
}
