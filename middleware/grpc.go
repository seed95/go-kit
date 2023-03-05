package middleware

import (
	"context"
	"fmt"
	"github.com/seed95/go-kit/log"
	"github.com/seed95/go-kit/log/keyval"
	"github.com/seed95/go-kit/log/zap"
	"google.golang.org/grpc"
	"time"
)

// LogInterceptor log request response
func LogInterceptor(l zap.Logger) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		start := time.Now()
		resp, err = handler(ctx, req)

		commonKeyVal := []keyval.Pair{
			keyval.String("req", fmt.Sprintf("%+v", req)),
			keyval.String("res", fmt.Sprintf("%+v", resp)),
		}
		log.ReqResWithLogger(l, start, err, commonKeyVal...)

		return resp, err

	}
}

// TimeoutInterceptor add deadline to context in request
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}

// RecoverInterceptor add deadline to context in request
func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Recover panic
		defer func() {
			if r := recover(); r != nil {
				resp = nil
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}
