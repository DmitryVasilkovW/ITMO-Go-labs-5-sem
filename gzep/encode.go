//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var pool sync.Pool

func Encode(data []byte, dst io.Writer) error {
	ww := getWriter(dst)

	defer func() {
		closeAndPut(ww)
	}()

	if _, err := ww.Write(data); err != nil {
		return err
	}

	return ww.Flush()
}

func getWriter(dst io.Writer) *gzip.Writer {
	var ww *gzip.Writer
	if w := pool.Get(); w != nil {
		ww = w.(*gzip.Writer)
		ww.Reset(dst)
	} else {
		ww, _ = gzip.NewWriterLevel(dst, gzip.DefaultCompression)
	}

	return ww
}

func closeAndPut(ww *gzip.Writer) {
	_ = ww.Close()
	pool.Put(ww)
}
