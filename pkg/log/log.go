package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TagTraceID   = "trace-id"
	TagSessionID = "session-id"
)

const (
	ContextKey = "ctx-logger"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05Z"))
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}

func Ctx(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(ContextKey).(*zap.SugaredLogger)
	if !ok {
		logger.Warn("Ctx doesn't contain logger")
		return logger
	}
	return l
}

func Global() *zap.SugaredLogger {
	return logger
}

// WithMeta sets logger meta tags
func WithMeta(meta map[string]any) []any {
	args := make([]any, 0, 2*len(meta))
	for k, v := range meta {
		args = append(args, k, v)
	}
	return args
}

// WrapCtx wraps the current logger in a context
func WrapCtx(ctx context.Context, meta map[string]any) context.Context {
	return context.WithValue(ctx, ContextKey, logger.With(WithMeta(meta)...))
}
