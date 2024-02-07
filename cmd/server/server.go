package main

import (
	"context"
	"fmt"
	stdlog "log"
	stdhttp "net/http"
	"time"

	"github.com/cameront/go-svelte-sqlite-template/config"
	"github.com/cameront/go-svelte-sqlite-template/http"
	"github.com/cameront/go-svelte-sqlite-template/logging"

	"github.com/cameront/go-svelte-sqlite-template/rpc_internal/counterservice"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, err := config.InitConfig(context.Background())
	panicIf(err, "error initializing config")
	cfg := config.MustContext(ctx)

	logger := logging.InitLogging(cfg)

	// receiving a signal on this channel keeps the server alive for another
	// IdleTimeoutMS
	requestsReceived := make(chan struct{}, 5)

	rpcHandler, apiClose, err := counterservice.InitApi(ctx, cfg.RPC_PathPrefix, requestsReceived)
	panicIf(err, "initializing api")

	staticHandler := http.InitStatic(ctx, cfg.HTTP_StaticDir)

	router, err := http.InitRouter(ctx, staticHandler, rpcHandler)
	panicIf(err, "initializing http")

	addr := fmt.Sprintf("%s:%s", cfg.RPC_Host, cfg.RPC_Port)
	httpServer := &stdhttp.Server{Addr: addr, Handler: router}
	closeServer := func() error {
		err := httpServer.Close()
		if err != stdhttp.ErrServerClosed {
			return err
		}
		return apiClose()
	}

	go startCloseTimer(ctx, cfg.HTTP_IdleShutdownMS, requestsReceived, []func() error{closeServer})

	logger.Info(fmt.Sprintf("listening on %s", addr))
	err = httpServer.ListenAndServe()
	if err != stdhttp.ErrServerClosed {
		panicIf(err, "server")
	}
}

// startCloseTimer calls all the closeFns provided if nothing is read from the
// requests channel within shutdownMS. This is helpful on fly.io's firecracker
// VMs so that we're only paying for used server time.
func startCloseTimer(ctx context.Context, shutdownMS int64, requests <-chan struct{}, closeFns []func() error) {
	if shutdownMS <= 0 {
		shutdownMS = 3000000000000 // ~100 years
	}
	duration := time.Duration(shutdownMS) * time.Millisecond

	t := time.NewTicker(time.Duration(duration))
	for {
		select {
		case <-requests:
			t.Reset(duration)
		case <-t.C:
			logging.GetLogger(ctx).Info(fmt.Sprintf("shutting down after %d ms", shutdownMS))
			for i, fn := range closeFns {
				if err := fn(); err != nil {
					logging.GetLogger(ctx).Error("error calling close fn %d: %v", i, err)
				}
			}
			return
		}
	}
}

func panicIf(err error, reason string) {
	if err != nil {
		stdlog.Fatalf("%s: %v", reason, err)
	}
}
