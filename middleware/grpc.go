package handler

import (
	"context"
	"fmt"
	"github.com/seed95/go-kit/log/keyval"
	"github.com/seed95/go-kit/log/zap"
	"google.golang.org/grpc"
	"time"
)

// TODO add comment
func LogInterceptor(l zap.Logger) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		start := time.Now()
		resp, err = handler(ctx, req)

		commonKeyVal := []keyval.Pair{
			keyval.String("duration", time.Since(start).String()),
			keyval.String("req", fmt.Sprintf("%+v", req)),
			keyval.String("res", fmt.Sprintf("%+v", resp)),
		}
		_ = commonKeyVal
		// TODO add logger as input
		//log.ReqRes(l, "grpc."+strings.Split(info.FullMethod, "/")[2], err, commonKeyVal...)

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

// TODO add comment
func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Recover panic
		defer func() {
			if r := recover(); r != nil {
				resp = nil
				//err = derror.New(derror.InternalServer, fmt.Sprintf("%+v", r))
			}
		}()

		return handler(ctx, req)
	}
}
