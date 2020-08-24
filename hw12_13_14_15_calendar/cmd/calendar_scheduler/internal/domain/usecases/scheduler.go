package usecases

import (
	"context"
)

type (
	Scheduler interface {
		SendAlerts(ctx context.Context)
		CleanEvents(ctx context.Context, daysToDelete int)
	}
)
