package usecases

import (
	"context"
)

type ILogger interface {
	Debug(ctx context.Context, message interface{}, args ...interface{})
	Info(ctx context.Context, message string, args ...interface{})
	Warn(ctx context.Context, message string, args ...interface{})
	Error(ctx context.Context, message interface{}, args ...interface{})
}
