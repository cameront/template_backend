package user

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cameront/template_backend/auth"
	"github.com/cameront/template_backend/ent"
	"github.com/cameront/template_backend/logging"
	rpc "github.com/cameront/template_backend/rpc/user"
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
	rpcHandler := rpc.NewUserServer(server,
		twirp.WithServerInterceptors(
			internal.NewLoggingInterceptor()),
		twirp.WithServerPathPrefix(pathPrefix))

	internal.HandleTwitchRPCAtPrefix(mux, pathPrefix, auth.UserAuthenticatingHandler(rpcHandler))

	return nil
}

func (s *Server) WhoAmI(ctx context.Context, req *rpc.Empty) (*rpc.Identity, error) {
	if claims, ok := ctx.Value(auth.UserCtxKey).(*auth.UserClaims); ok {
		return &rpc.Identity{
			Username: claims.Name,
			IsAdmin:  false,
		}, nil
	} else {
		logging.Errorf(ctx, "how did a user bypass the authenticating handler!?")
		return nil, twirp.InternalError("uh oh")
	}
}

func (s *Server) Logout(ctx context.Context, req *rpc.Empty) (*rpc.Empty, error) {
	cookie := http.Cookie{
		Name:    auth.AuthCookieName,
		Expires: time.Unix(0, 0),
		Value:   "",
		Path:    "/",
		MaxAge:  0,
	}
	twirp.AddHTTPResponseHeader(ctx, "Set-Cookie", cookie.String())

	return &rpc.Empty{}, nil
}

func (s *Server) GetValue(ctx context.Context, req *rpc.CounterRequest) (*rpc.CounterValue, error) {

	op := func(ctx context.Context, db *ent.Client) (int64, error) {
		c, err := db.Counter.Get(ctx, strings.ToLower(req.Name))

		if err == nil {
			return c.Value, nil
		}

		if ent.IsNotFound(err) {
			_, err = createCounter(ctx, db, req.Name, 0)
			return 0, err
		}

		return 0, err
	}

	val, err := store.WithTransaction(ctx, s.db, op)
	if err != nil {
		return nil, err
	}

	return &rpc.CounterValue{Name: req.Name, Value: val}, nil
}

func (s *Server) Increment(ctx context.Context, req *rpc.IncrementRequest) (*rpc.CounterValue, error) {

	op := func(ctx context.Context, db *ent.Client) (int64, error) {
		c, err := db.Counter.Get(ctx, strings.ToLower(req.Name))
		if err != nil {
			if !ent.IsNotFound(err) {
				return 0, err
			}
			c, err = createCounter(ctx, db, req.Name, 0)
			if err != nil {
				return 0, err
			}
		}

		update := c.Update()
		mut := update.Mutation()
		newValue := c.Value + 1
		mut.SetValue(newValue)
		_, err = update.Save(ctx)

		return newValue, err
	}

	val, err := store.WithTransaction(ctx, s.db, op)
	if err != nil {
		return nil, err
	}

	return &rpc.CounterValue{Name: req.Name, Value: val}, nil
}

func createCounter(ctx context.Context, db *ent.Client, name string, value int64) (*ent.Counter, error) {
	create := db.Counter.Create()
	mut := create.Mutation()
	mut.SetID(strings.ToLower(name))
	mut.SetValue(value)
	return create.Save(ctx)
}
