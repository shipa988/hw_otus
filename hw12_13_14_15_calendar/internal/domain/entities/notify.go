package entities

import (
	"context"
	"time"
)

type NotifyRepo interface {
	Add(ctx context.Context, alert Notify) error
}

type Notify struct {
	ID       string
	Title    string
	UserID   string
	DateTime time.Time
}
