package rpcinternal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cameront/template_backend/auth"
	"github.com/cameront/template_backend/logging"

	"github.com/segmentio/ksuid"
	"github.com/twitchtv/twirp"
)

// Logs RPC methods on the way in with payload, and the way out with
// duration/result.
func NewLoggingInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			start := time.Now()

			user := ""
			if claims, ok := ctx.Value(auth.UserCtxKey).(*auth.UserClaims); ok {
				user = claims.Name
			}

			serviceName, _ := twirp.ServiceName(ctx)
			methodName, _ := twirp.MethodName(ctx)
			traceLogger := logging.GetLogger(ctx).
				With("traceId", ksuid.New()).
				With("user", user)
			ctx = logging.SetLogger(ctx, traceLogger)

			traceLogger.
				With("initiator", "rpc").
				With("method", fmt.Sprintf("%s.%s", serviceName, methodName)).
				With("payload", fmt.Sprintf("%+v", req)).
				Info("begin request")
			res, err := next(ctx, req)
			traceLogger.
				With("durationms", time.Since(start).Milliseconds()).
				With("error", err).
				Info("completed request")

			return res, err
		}
	}
}

func HandleTwitchRPCAtPrefix(mux *http.ServeMux, prefix string, handler http.Handler) error {
	if !strings.HasPrefix(prefix, "/") {
		return fmt.Errorf("path prefix (%s) must begin with '/'", prefix)
	}
	if strings.HasSuffix(prefix, "/") {
		return fmt.Errorf("path prefix (%s) must not end with '/'", prefix)
	}

	path := fmt.Sprintf("POST %s/", prefix)
	mux.Handle(path, handler)

	return nil
}
