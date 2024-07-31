package handlers

import (
	"chi_boilerplate/pkg/infrastructure/logger"
	"net/http"
)

// WrapError wraps the handlers error and logs it
func WrapError(f func(w http.ResponseWriter, r *http.Request) error, l *logger.ZapLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			l.Error(err.Error())
		}
	}
}
