package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/seed95/go-kit/log"
	"github.com/seed95/go-kit/log/keyval"
	"github.com/seed95/go-kit/log/zap"
)

// LogMiddleware Log request
func LogMiddleware(l zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info("http_request", keyval.String("req", fmt.Sprintf("%+v", r)))
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

// TimeoutMiddleware add deadline to context in request
func TimeoutMiddleware(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
		next(w, r)
	}
}
