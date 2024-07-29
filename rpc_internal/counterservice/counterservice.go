package counterservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cameront/template_backend/auth"
	"github.com/cameront/template_backend/ent"
	rpc "github.com/cameront/template_backend/rpc/count"
	internal "github.com/cameront/template_backend/rpc_internal"
	"github.com/cameront/template_backend/store"

	"github.com/twitchtv/twirp"
)

// Server implements the Counter service
type Server struct {
	db store.DbProvider
}

func InitApi(ctx context.Context, dbProvider store.DbProvider, mux *http.ServeMux, pathPrefix string, requestsReceived chan<- struct{}) (func() error, error) {

	server := &Server{db: dbProvider}
	rpcHandler := rpc.NewCounterServer(server,
		twirp.WithServerInterceptors(
			internal.NewSignalInterceptor(requestsReceived),
			internal.NewLoggingInterceptor()),
		twirp.WithServerPathPrefix(pathPrefix))

	if !strings.HasPrefix(pathPrefix, "/") {
		return nil, fmt.Errorf("path prefix (%s) must begin with '/'", pathPrefix)
	}
	if strings.HasSuffix(pathPrefix, "/") {
		return nil, fmt.Errorf("path prefix (%s) must not end with '/'", pathPrefix)
	}

	path := fmt.Sprintf("POST %s/", pathPrefix)
	mux.Handle(path, auth.UserAuthenticatingHandler(rpcHandler))

	closeFn := func() error {
		c := dbProvider.Get(ctx)
		c.Close()
		return nil
	}
	return closeFn, nil
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
