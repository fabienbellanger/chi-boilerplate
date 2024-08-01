package handlers

import (
	"chi_boilerplate/pkg/infrastructure/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap/zapcore"
)

// RequestIDKey is the key used to store the request ID in the context
type RequestIDKey string

// WrapError wraps the handlers error and logs it
func WrapError(f func(w http.ResponseWriter, r *http.Request) error, l *logger.ZapLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := fmt.Sprintf("%s", r.Context().Value(RequestIDKey("request_id")))

		err := f(w, r)
		if err != nil {
			l.Error(err.Error(), zapcore.Field{Key: "request_id", Type: zapcore.StringType, String: requestId})
		}
	}
}
