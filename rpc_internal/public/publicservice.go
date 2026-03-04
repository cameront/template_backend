package public

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cameront/template_backend/auth"
	rpc "github.com/cameront/template_backend/rpc/public"
	internal "github.com/cameront/template_backend/rpc_internal"
	"github.com/cameront/template_backend/store"

	"github.com/twitchtv/twirp"
)

// Server implements the Counter service
type Server struct {
	db store.DbProvider
}

func InitApi(ctx context.Context, dbProvider store.DbProvider, mux *http.ServeMux, pathPrefix string) error {

	server := &Server{db: dbProvider}
	rpcHandler := rpc.NewPublicServer(server,
		twirp.WithServerInterceptors(
			internal.NewLoggingInterceptor()),
		twirp.WithServerPathPrefix(pathPrefix))

	internal.HandleTwitchRPCAtPrefix(mux, pathPrefix, rpcHandler)

	return nil
}

func (s *Server) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.IdentityMessage, error) {

	if req.Username != "meuser" || req.Password != "pass123" {
		return nil, twirp.Unauthenticated.Error("invalid username/password")
	}

	twoDays := time.Hour * 24 * 2
	expires := time.Now().Add(twoDays)
	token, err := auth.BuildToken("123", req.Username, req.Username+"@example.com", "none", expires)
	if err != nil {
		return nil, twirp.InternalError(fmt.Sprintf("building token: %w", err))
	}

	cookie := http.Cookie{Name: auth.AuthCookieName, Value: token, Expires: expires, Path: "/"}
	twirp.AddHTTPResponseHeader(ctx, "Set-Cookie", cookie.String())

	return &rpc.IdentityMessage{Authenticated: true, Username: req.Username, IsAdmin: false}, nil
}
