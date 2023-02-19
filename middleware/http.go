package middleware

import (
	"fmt"
	"github.com/seed95/go-kit/log"
	"github.com/seed95/go-kit/log/keyval"
	"net/http"
	"runtime/debug"
)

// LogMiddleware Log request
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("http_request", keyval.String("req", fmt.Sprintf("%+v", r)))
		next(w, r)
	}
}

// RecoverMiddleware recover panic as return error
func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Recover panic
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				log.Error("panic", keyval.String("stacktrace", stack))
			}
		}()

		next(w, r)
	}
}
