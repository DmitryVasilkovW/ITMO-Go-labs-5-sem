package otp

import "io"

type cipherStreamReader struct {
	r    io.Reader
	prng io.Reader
}

func (cr *cipherStreamReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	if n > 0 {
		prngData := make([]byte, n)
		_, _ = cr.prng.Read(prngData)

		for i := 0; i < n; i++ {
			p[i] ^= prngData[i]
		}
	}

	return n, err
}
