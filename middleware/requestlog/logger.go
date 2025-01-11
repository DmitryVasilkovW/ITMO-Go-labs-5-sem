package requestlog

import (
	"time"

	"go.uber.org/zap"
	"net/http"
)

func logRequestStart(l *zap.Logger, requestID string, r *http.Request) {
	l.Info("request started",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)
}

func logRequestFinish(l *zap.Logger, requestID string, duration time.Duration, wrapper *ResponseWriter, r *http.Request) {
	l.Info("request finished",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.Duration("duration", duration),
		zap.Int("status_code", wrapper.status),
	)
}

func logPanic(l *zap.Logger, requestID string, duration time.Duration, wrapper *ResponseWriter, r *http.Request, err interface{}) {
	stackTrace := captureStackTrace()
	l.Error("request panicked",
		zap.String("request_id", requestID),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.Duration("duration", duration),
		zap.Int("status_code", wrapper.status),
		zap.String("stack_trace", stackTrace),
		zap.Any("panic_error", err),
	)
}
