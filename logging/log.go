package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/cameront/go-svelte-sqlite-template/config"
)

type contextKeyType string

const ctxKey contextKeyType = "_log"

var log *slog.Logger

func InitLogging(cfg *config.Config) *slog.Logger {
	handlerOptions := &slog.HandlerOptions{AddSource: true}
	var handler slog.Handler
	if cfg.LOG_OutputFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}

	log = slog.New(handler)
	return log
}

func With(parent *slog.Logger, key, val interface{}) *slog.Logger {
	if parent == nil {
		parent = log
	}
	return parent.With(key, val)
}

func GetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxKey).(*slog.Logger)
	if !ok {
		return log
	}
	return logger
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey, logger)
}
