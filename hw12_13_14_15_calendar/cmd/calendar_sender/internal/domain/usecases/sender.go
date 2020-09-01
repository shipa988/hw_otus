package usecases

import (
	"context"
)

type (
	Sender interface {
		SendingAlerts(ctx context.Context) error
	}
)
