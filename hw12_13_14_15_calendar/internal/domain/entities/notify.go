package entities

import (
	"context"
	"time"
)

type NotifyQueue interface {
	Pull(ctx context.Context, alert chan<- Notify) error
	Push(Notify) error
}

type Notify struct {
	ID       string
	Title    string
	UserID   string
	DateTime time.Time
}
