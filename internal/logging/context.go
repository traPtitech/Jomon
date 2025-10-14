package logging

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type contextKey struct{}

var ErrLoggerNotSet = errors.New("logger not set in context")

func SetLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func GetLoggerMaybe(ctx context.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value(contextKey{}).(*zap.Logger)
	if !ok {
		return nil, ErrLoggerNotSet
	}
	return logger, nil
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, err := GetLoggerMaybe(ctx)
	if err != nil {
		panic(err)
	}
	return logger
}
