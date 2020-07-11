package mocks

import (
	"context"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

var _ usecases.Logger = (*MockLogger)(nil)

type MockLogger struct{}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m MockLogger) Debug(ctx context.Context, message interface{}, args ...interface{}) {
}

func (m MockLogger) Info(ctx context.Context, message string, args ...interface{}) {
}

func (m MockLogger) Warn(ctx context.Context, message string, args ...interface{}) {
}

func (m MockLogger) Error(ctx context.Context, message interface{}, args ...interface{}) {
}
