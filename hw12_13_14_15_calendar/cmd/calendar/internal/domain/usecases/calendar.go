package usecases

import (
	"context"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
)

type Calendar interface {
	MakeEvent(ctx context.Context, title, dateTimeEvent, text, userID, duration, timeNotify string) (id string, err error)
	UpdateEvent(ctx context.Context, userID, eventID, newTitle, newDateTimeEvent, newText, newDuration, newTimeNotify string) (id string, err error)
	DeleteEvent(ctx context.Context, userID, eventID string) (id string, err error)
	GetDateEvents(ctx context.Context, date, userID string) ([]*entities.Event, error)
	GetWeekEvents(ctx context.Context, date, userID string) ([]*entities.Event, error)
	GetMonthEvents(ctx context.Context, date, userID string) ([]*entities.Event, error)
}
