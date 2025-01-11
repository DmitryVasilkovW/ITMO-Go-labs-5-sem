package requestlog

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GenerateRequestID()
			start := time.Now()
			wrapper := &ResponseWriter{ResponseWriter: w, status: http.StatusOK}

			handleRequest(l, next, wrapper, r, requestID, start)
		})
	}
}

func handleRequest(l *zap.Logger, next http.Handler, wrapper *ResponseWriter, r *http.Request, requestID string, start time.Time) {
	logRequestStart(l, requestID, r)

	defer func() {
		duration := time.Since(start)

		if err := recover(); err != nil {
			logPanic(l, requestID, duration, wrapper, r, err)
			panic(err)
		}

		logRequestFinish(l, requestID, duration, wrapper, r)
	}()

	next.ServeHTTP(wrapper, r)
}
