package http

import (
	"context"
	"net/http"
	"path"

	"github.com/cameront/go-svelte-sqlite-template/auth"
	"github.com/cameront/go-svelte-sqlite-template/config"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func InitRouter(ctx context.Context, staticHandler http.Handler, rpcAPI http.Handler) (*httprouter.Router, error) {
	cfg := config.MustContext(ctx)

	corsWrapper := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"OPTIONS", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	router := httprouter.New()
	router.RedirectTrailingSlash = true

	// Handle (non-authenticated, obviously) login requests
	router.Handler("POST", "/login", corsWrapper.Handler(auth.LoginHandler()))

	// Handle options - though this handler is never called because the
	// corsWrapper doesn't pass through options requests, but we need a
	// handler to wrap... so.
	nopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	router.Handler("OPTIONS", "/*path", corsWrapper.Handler(nopHandler))

	// Handle health checks
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO: query something simple from the db to get a better "health" check
		w.Write([]byte("ok"))
	})

	// Handle RPC requests
	auth.InitJWT(cfg.AUTH_JWTSecret)
	namedPrefix := path.Join(cfg.RPC_PathPrefix, "/*func")
	router.Handler("POST", namedPrefix, corsWrapper.Handler(
		auth.UserAuthenticatingHandler(
			rpcAPI,
		),
	),
	)

	// A little hacky, but I'm not the first one that's thought of it, since
	// it's "recommended" on the httprouter docs.
	router.NotFound = staticHandler

	return router, nil
}
