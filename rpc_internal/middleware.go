package rpcinternal

import (
	"context"
	"fmt"
	"time"

	"github.com/cameront/go-svelte-sqlite-template/auth"
	"github.com/cameront/go-svelte-sqlite-template/log"
	"github.com/segmentio/ksuid"
	"github.com/twitchtv/twirp"
)

func NewSignalInterceptor(requestsSignal chan<- struct{}) twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			requestsSignal <- struct{}{}
			return next(ctx, req)
		}
	}
}

func NewLoggingInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			start := time.Now()

			serviceName, _ := twirp.ServiceName(ctx)
			methodName, _ := twirp.MethodName(ctx)
			logger := log.GetLogger(ctx).
				With("initiator", "rpc").
				With("traceId", ksuid.New()).
				With("method", fmt.Sprintf("%s.%s", serviceName, methodName)).
				With("user", fmt.Sprintf("%s", ctx.Value(auth.UserCtxKey))).
				With("payload", fmt.Sprintf("%+v", req))
			ctx = log.SetLogger(ctx, logger)

			logger.Info("begin request")
			res, err := next(ctx, req)
			resLogger := logger.With("durationms", time.Since(start).Milliseconds()).With("error", err)
			resLogger.Info("completed request")

			return res, err
		}
	}
}
