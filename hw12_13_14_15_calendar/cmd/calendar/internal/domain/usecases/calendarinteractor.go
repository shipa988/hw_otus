package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

var _ Calendar = (*CalendarInteractor)(nil)

const (
	ErrMake           = "can't make event in calendar"
	ErrUpdate         = "can't update event in calendar by id: %v"
	ErrDelete         = "can't delete event in calendar by id: %v"
	ErrGetDateEvents  = "can't get events for date %v from calendar"
	ErrGetWeekEvents  = "can't get events for week starts from %v from calendar"
	ErrGetMonthEvents = "can't get events for month starts from %v from calendar"
)

type CalendarInteractor struct {
	events entities.EventRepo
	alerts entities.NotifyQueue
	logger usecases.Logger
}

func NewCalendar(eRepo entities.EventRepo, nRepo entities.NotifyQueue, logger usecases.Logger) *CalendarInteractor {
	return &CalendarInteractor{
		events: eRepo,
		alerts: nRepo,
		logger: logger,
	}
}

func (c CalendarInteractor) MakeEvent(ctx context.Context, title, dateTimeEvent, text, userID, duration, timeNotify string) (string, error) {
	event, err := entities.NewEvent(title, dateTimeEvent, duration, text, userID, timeNotify)
	if err != nil {
		return "", errors.Wrap(err, ErrMake)
	}
	id, err := c.events.Add(ctx, *event)
	if err != nil {
		return "", errors.Wrap(err, ErrMake)
	}
	c.logger.Info(ctx, "Event id: %v added to calendar", id)
	return id, nil
}

func (c CalendarInteractor) UpdateEvent(ctx context.Context, userID, eventID, newTitle, newDateTimeEvent, newText, newDuration, newTimeNotify string) (string, error) {
	if eventID == "" {
		return "", errors.Wrap(fmt.Errorf(entities.ErrNoField, "id"), ErrMake)
	}
	if userID == "" {
		return "", errors.Wrap(fmt.Errorf(entities.ErrNoField, "userid"), ErrMake)
	}

	e, err := c.events.GetByID(ctx, userID, eventID)
	if err != nil {
		return "", errors.Wrapf(err, ErrUpdate, eventID)
	}

	upde, err := e.Update(newTitle, newDateTimeEvent, newDuration, newText, newTimeNotify)
	if err != nil {
		return "", errors.Wrapf(err, ErrUpdate, eventID)
	}
	if err := c.events.UpdateByID(ctx, userID, eventID, upde); err != nil {
		return "", errors.Wrapf(err, ErrUpdate, eventID)
	}
	c.logger.Info(ctx, "Event id: %v updated in calendar", eventID)
	return eventID, nil
}

func (c CalendarInteractor) DeleteEvent(ctx context.Context, userID, eventID string) (string, error) {
	if eventID == "" {
		return "", errors.Wrap(fmt.Errorf(entities.ErrNoField, "id"), ErrMake)
	}
	if userID == "" {
		return "", errors.Wrap(fmt.Errorf(entities.ErrNoField, "userid"), ErrMake)
	}

	if err := c.events.DeleteByUserID(ctx, userID, eventID); err != nil {
		return "", errors.Wrapf(err, ErrDelete, eventID)
	}
	c.logger.Info(ctx, "Event id: %v deleted from calendar", eventID)
	return eventID, nil
}

func (c CalendarInteractor) GetDateEvents(ctx context.Context, date, userID string) ([]*entities.Event, error) {
	if date == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "date"), ErrMake)
	}
	if userID == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "userid"), ErrMake)
	}
	dt, err := time.Parse(entities.LayoutDateISO, date)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrDateFormat, "error parse dateTime:  %v", date)
	}
	events, err := c.events.GetByDate(ctx, userID, dt)
	if err != nil {
		return nil, errors.Wrapf(err, ErrGetDateEvents, date)
	}
	return events, nil
}

func (c CalendarInteractor) GetWeekEvents(ctx context.Context, date, userID string) ([]*entities.Event, error) {
	if date == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "date"), ErrMake)
	}
	if userID == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "userid"), ErrMake)
	}
	dt, err := time.Parse(entities.LayoutDateISO, date)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrDateFormat, "error parse dateTime:  %v", date)
	}
	events, err := c.events.GetForPeriodByUserID(ctx, userID, dt, dt.AddDate(0, 0, 7))
	if err != nil {
		return nil, errors.Wrapf(err, ErrGetWeekEvents, date)
	}
	return events, nil
}

func (c CalendarInteractor) GetMonthEvents(ctx context.Context, date, userID string) ([]*entities.Event, error) {
	if date == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "date"), ErrMake)
	}
	if userID == "" {
		return nil, errors.Wrap(fmt.Errorf(entities.ErrNoField, "userid"), ErrMake)
	}
	dt, err := time.Parse(entities.LayoutDateISO, date)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrDateFormat, "error parse dateTime:  %v", date)
	}
	events, err := c.events.GetForPeriodByUserID(ctx, userID, dt, dt.AddDate(0, 1, 0))
	if err != nil {
		return nil, errors.Wrapf(err, ErrGetMonthEvents, date)
	}
	return events, nil
}
