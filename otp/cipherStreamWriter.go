package otp

import "io"

type cipherStreamWriter struct {
	w    io.Writer
	prng io.Reader
}

func (cw *cipherStreamWriter) Write(p []byte) (int, error) {
	bytes := make([]byte, len(p))
	prngData := make([]byte, len(p))
	_, _ = cw.prng.Read(prngData)

	for i := 0; i < len(p); i++ {
		bytes[i] = p[i] ^ prngData[i]
	}

	return cw.w.Write(bytes)
}
