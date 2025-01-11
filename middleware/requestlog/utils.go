package requestlog

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"time"
)

func GenerateRequestID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return time.Now().Format(time.RFC3339Nano)
	}
	return hex.EncodeToString(bytes)
}

func captureStackTrace() string {
	bytes := make([]byte, 4096)
	bytes = bytes[:runtime.Stack(bytes, false)]
	return string(bytes)
}
