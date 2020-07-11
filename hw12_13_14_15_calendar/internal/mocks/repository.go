package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
)

var _ entities.EventRepo = (*EventRepo)(nil)

type MapRepo struct {
}

type EventRepo struct {
	testEvent *entities.Event
}

func NewMockRepo(event *entities.Event) *EventRepo {
	return &EventRepo{
		testEvent: event,
	}
}

func (i *EventRepo) Add(ctx context.Context, event entities.Event) (string, error) {
	return i.testEvent.ID, nil
}

func (i *EventRepo) GetByID(ctx context.Context, userID, eventID string) (*entities.Event, error) {
	return i.testEvent, nil
}

func (i *EventRepo) GetByDate(ctx context.Context, userID string, date time.Time) ([]*entities.Event, error) {
	return []*entities.Event{i.testEvent}, nil
}

func (i *EventRepo) GetForPeriod(ctx context.Context, userID string, dateStart time.Time, dateEnd time.Time) ([]*entities.Event, error) {
	return []*entities.Event{i.testEvent}, nil
}

func (i *EventRepo) UpdateByID(ctx context.Context, userID, eventID string, event entities.Event) error {
	if userID == "" || eventID == "" {
		return errors.New("")
	}
	return nil
}

func (i *EventRepo) DeleteByID(ctx context.Context, userID, eventID string) error {
	if userID == "" || eventID == "" {
		return errors.New("")
	}
	return nil
}
