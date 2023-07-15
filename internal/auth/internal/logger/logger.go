package logger

// In the logger package, we use the nesting of the logger in the context,
// for a more convenient transfer of the logger through the context.
//
// InitLog provides us with a function to get a context with a logger, and the logger is initialized once thanks to the sync.Once
// For use this package call InitLog() passing log level, after passing

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogger struct{}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func InitLogger(level string) (*zap.Logger, error) {

	var lvl zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)

	// lvl, err = zap.ParseAtomicLevel("Debug")

	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	log := zl
	log.Info(`Logger level`, zap.String("logLevel", level))

	return log, nil
}
