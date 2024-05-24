package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type contextKeyType string

const ctxKey contextKeyType = "_log"

var log *slog.Logger

func InitLogging(outputFormat string, minLogLevel int) *slog.Logger {
	handlerOptions := &slog.HandlerOptions{AddSource: true, Level: slog.Level(minLogLevel)}
	var handler slog.Handler
	if outputFormat == "json" {
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

func Debugf(ctx context.Context, format string, args ...any) {
	logSkipCaller(ctx, GetLogger(ctx), slog.LevelDebug, format, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	logSkipCaller(ctx, GetLogger(ctx), slog.LevelInfo, format, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	logSkipCaller(ctx, GetLogger(ctx), slog.LevelWarn, format, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	logSkipCaller(ctx, GetLogger(ctx), slog.LevelError, format, args...)
}

func logSkipCaller(ctx context.Context, logger *slog.Logger, level slog.Level, format string, args ...any) {
	if !logger.Enabled(ctx, level) {
		return
	}
	var pcs [2]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[1])
	_ = logger.Handler().Handle(ctx, r)
}
