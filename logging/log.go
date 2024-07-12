package logging

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type contextKeyType string

const ctxKey contextKeyType = "_log"

var log *slog.Logger

// If we can tell the package path, we strip that from our log lines because
// it's a lot of text that is all very repetitive "github.com/foo/bar/zee/..."
var packagePath = ""

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		stdlog.Println("no debug.ReadBuildInfo available. perhaps not building in module mode!?")
		return
	}
	if info.Path == "command-line-arguments" {
		// This usually indicates that you're building in a way that doesn't
		// include the module info. Perhaps you're building like:
		// go build -o /some/bin ./cmd/app/main.go
		// instead of go build -o /some/bin ./cmd/app?
		stdlog.Println("module appears misconfigured (info.Path = 'command-line-arguments')")
		return
	}
	packagePath = info.Main.Path
}

func InitLogging(outputFormat string, minLogLevel int) *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
				if packagePath != "" {
					s.Function = strings.TrimPrefix(s.Function, packagePath)
				}
			}
			return a
		},
		Level: slog.Level(minLogLevel)}
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
