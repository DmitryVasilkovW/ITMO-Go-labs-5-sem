//go:build !solution

package otp

import (
	"io"
)

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &cipherStreamReader{r: r, prng: prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &cipherStreamWriter{w: w, prng: prng}
}
