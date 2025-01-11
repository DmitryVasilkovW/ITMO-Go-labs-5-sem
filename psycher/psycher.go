//go:build !solution

package psycher

type Cipher struct {
	blockSize int
	keys      [128][16]byte
}

func New(keys [][]byte) *Cipher {
	var c Cipher
	c.blockSize = 16

	for i, key := range keys {
		copy(c.keys[i][:], key)
	}

	return &c
}

func (c Cipher) BlockSize() int {
	return c.blockSize
}

func (c Cipher) Encrypt(dst, src []byte) {
	for i := range dst {
		dst[i] = 0
	}

	for bitIndex := 0; bitIndex < 128; bitIndex++ {
		if (src[bitIndex/8] & (1 << (bitIndex % 8))) != 0 {
			for j := 0; j < c.blockSize; j++ {
				dst[j] ^= c.keys[bitIndex][j]
			}
		}
	}
}

func (c Cipher) Decrypt(dst, src []byte) {
	panic("Decrypt not implemented")
}
