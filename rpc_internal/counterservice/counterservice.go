package counterservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cameront/go-svelte-sqlite-template/ent"
	rpc "github.com/cameront/go-svelte-sqlite-template/rpc/count"
	internal "github.com/cameront/go-svelte-sqlite-template/rpc_internal"

	"github.com/twitchtv/twirp"
)

// Server implements the Counter service
type Server struct {
	DbClient *clientWrapper
}

func InitApi(ctx context.Context, pathPrefix string, requestsReceived chan<- struct{}) (http.Handler, func() error, error) {
	// We assume migrations have already been applied by atlas in this point
	// by either the dev environment (air script) or in production (docker
	// step)
	dbClient, err := InitStore(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing database %v", err)
	}

	server := &Server{DbClient: dbClient} // implements service interface
	rpcHandler := rpc.NewCounterServer(server,
		twirp.WithServerInterceptors(
			internal.NewSignalInterceptor(requestsReceived),
			internal.NewLoggingInterceptor()),
		twirp.WithServerPathPrefix(pathPrefix))

	closeFn := func() error {
		dbClient.Close()
		return nil
	}
	return rpcHandler, closeFn, nil
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

	val, err := WithTransaction(ctx, s.DbClient, op)
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

	val, err := WithTransaction(ctx, s.DbClient, op)
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
