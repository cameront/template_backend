package main

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cameront/template_backend/config"
	"github.com/cameront/template_backend/logging"
	"github.com/cameront/template_backend/static"
	"github.com/cameront/template_backend/store"
	"github.com/rs/cors"

	rpcpublic "github.com/cameront/template_backend/rpc_internal/public"
	rpcuser "github.com/cameront/template_backend/rpc_internal/user"

	_ "github.com/mattn/go-sqlite3"
)

var isShuttingDown atomic.Bool

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx, err := config.InitConfig(ctx)
	panicIf(err, "initializing config")
	cfg := config.MustContext(ctx)

	logging.InitLogging(cfg.LOG_OutputFormat, cfg.LOG_MinLevel)

	dbClient, err := store.InitStore(ctx)
	panicIf(err, "initializind db")

	mux := http.NewServeMux()

	serveHealth(mux)

	err = rpcpublic.InitApi(ctx, dbClient, mux, "/rpc/public")
	panicIf(err, "initializing public api")

	err = rpcuser.InitApi(ctx, dbClient, mux, "/rpc/user")
	panicIf(err, "initializing user api")

	static.InitStatic(ctx, mux, cfg.HTTP_StaticDir)

	c := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{http.MethodPost, http.MethodOptions},
		AllowCredentials: true,
	})

	addr := fmt.Sprintf("%s:%s", cfg.RPC_Host, cfg.RPC_Port)
	httpServer := &http.Server{
		Addr:        addr,
		BaseContext: func(net.Listener) context.Context { return ctx },
		Handler:     c.Handler(mux),
	}

	runServer(ctx, httpServer)
}

func runServer(ctx context.Context, server *http.Server) {
	go func() {
		logging.Infof(ctx, "server listening at %s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panicIf(err, "server listening")
		}
		logging.Infof(ctx, "stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	isShuttingDown.Store(true)

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	// somewhere here we should be releasing resources

	if err := server.Shutdown(shutdownCtx); err != nil {
		server.Close()
		panicIf(err, "shutting down server")
	}

	logging.Infof(ctx, "shutdown completed gracefully")
}

func serveHealth(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "shutting down", http.StatusServiceUnavailable)
			return
		}
		time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "OK")
	})
}

func panicIf(err error, reason string) {
	if err != nil {
		stdlog.Fatalf("%s: %v", reason, err)
	}
}
