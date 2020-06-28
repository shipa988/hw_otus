package usecases

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/interfaces"
)

const (
	LayoutISO = "2006-01-02 15:04:05"
)

var _ interfaces.Calendar = (*CalendarInteractor)(nil)

type CalendarInteractor struct {
	events domain.EventRepo
	alerts domain.NotifyRepo
	logger ILogger
}

func NewCalendar(eRepo domain.EventRepo, nRepo domain.NotifyRepo, logger ILogger) *CalendarInteractor {
	return &CalendarInteractor{
		events: eRepo,
		alerts: nRepo,
		logger: logger,
	}
}

func (c CalendarInteractor) MakeEvent(ctx context.Context, title, dateTimeEvent, text, userID string, duration, timeNotify time.Duration) error {
	dt, err := time.Parse(LayoutISO, dateTimeEvent)
	if err != nil {
		return errors.Wrap(err, "can't make event in calendar\n")
	}
	event := domain.Event{
		Title:      title,
		DateTime:   dt,
		Duration:   duration,
		Text:       text,
		UserID:     userID,
		TimeNotify: timeNotify,
	}
	if err = event.Validate(); err != nil {
		c.logger.Warn(ctx, "Validation error %v", err)
	}
	id, err := c.events.Add(ctx, event)
	if err != nil {
		return errors.Wrap(err, "can't make event in calendar\n")
	}
	c.logger.Info(ctx, "Event id:%v add to calendar", id)
	return nil
}

func (c CalendarInteractor) UpdateEvent(ctx context.Context, id string, newTitle, newDateTimeEvent, newText, newUserID string, newDuration, newTimeNotify time.Duration) error {
	dt, err := time.Parse(LayoutISO, newDateTimeEvent)
	if err != nil {
		return errors.Wrap(err, "can't make event in calendar\n")
	}
	event := domain.Event{
		Title:      newTitle,
		DateTime:   dt,
		Duration:   newDuration,
		Text:       newText,
		UserID:     newUserID,
		TimeNotify: newTimeNotify,
	}

	if err := c.events.UpdateByID(ctx, id, event); err != nil {
		return errors.Wrap(err, "can't make event in calendar\n")
	}
	c.logger.Info(ctx, "Event %v update to calendar", newTitle)
	return nil
}

func (c CalendarInteractor) DeleteEvent(ctx context.Context, id string) error {
	if err := c.events.DeleteByID(ctx, id); err != nil {
		return errors.Wrap(err, "can't delete event from calendar\n")
	}
	return nil
}

func (c CalendarInteractor) GetDateEvents(ctx context.Context, date time.Time) ([]*domain.Event, error) {
	events, err := c.events.GetByDate(ctx, date)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get events for date %v from calendar\n", date.String())
	}
	return events, nil
}

func (c CalendarInteractor) GetWeekEvents(ctx context.Context, date time.Time) ([]*domain.Event, error) {
	events, err := c.events.GetForPeriod(ctx, date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, errors.Wrapf(err, "can't get events for week starts from %v from calendar\n", date.String())
	}
	return events, nil
}

func (c CalendarInteractor) GetMonthEvents(ctx context.Context, date time.Time) ([]*domain.Event, error) {
	events, err := c.events.GetForPeriod(ctx, date, date.AddDate(0, 1, 0))
	if err != nil {
		return nil, errors.Wrapf(err, "can't get events for month starts from %v from calendar\n", date.String())
	}
	return events, nil
}
