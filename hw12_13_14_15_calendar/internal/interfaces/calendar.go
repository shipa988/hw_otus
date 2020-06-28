package interfaces

import (
	"context"
	"time"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain"
)

type Calendar interface {
	MakeEvent(ctx context.Context, title, dateTimeEvent, text, userID string, duration, timeNotify time.Duration) error
	UpdateEvent(ctx context.Context, id string, newTitle, newDateTimeEvent, newText, newUserID string, newDuration, newTimeNotify time.Duration) error
	DeleteEvent(ctx context.Context, id string) error
	GetDateEvents(ctx context.Context, date time.Time) ([]*domain.Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]*domain.Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]*domain.Event, error)
}
