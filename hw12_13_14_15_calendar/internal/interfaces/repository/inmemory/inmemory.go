package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/usecases"
)

var _ domain.EventRepo = (*EventRepo)(nil)

type MapRepo struct {
	rwmux  *sync.RWMutex
	dates  map[time.Time]string
	events map[string]*domain.Event
}

func NewMapRepo() *MapRepo {
	return &MapRepo{
		rwmux:  &sync.RWMutex{},
		dates:  make(map[time.Time]string),
		events: make(map[string]*domain.Event),
	}
}
func (m *MapRepo) Clear() {
	m.rwmux.Lock()
	defer m.rwmux.Unlock()
	m.events = make(map[string]*domain.Event)
	m.dates = make(map[time.Time]string)
}

type EventRepo struct {
	m      *MapRepo //emulate infrastructure layer over type MapRepo
	logger usecases.ILogger
}

func NewInMemoryEventRepo(m *MapRepo, logger usecases.ILogger) (*EventRepo, error) {
	return &EventRepo{
		m:      m,
		logger: logger,
	}, nil
}

func (i EventRepo) Add(ctx context.Context, event domain.Event) (string, error) {
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()
	if _, ok := i.m.dates[event.DateTime]; ok {
		return "", domain.ErrDateBusy
	}
	id := uuid.NewV4().String()
	event.ID = id
	i.m.events[id] = &event
	i.m.dates[event.DateTime] = event.ID
	return id, nil
}

func (i EventRepo) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	e, ok := i.m.events[id]
	if !ok {
		return nil, domain.ErrEventNotFound
	}
	return eventCreateSafely(e), nil
}

func (i EventRepo) GetByDate(ctx context.Context, date time.Time) ([]*domain.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	events := []*domain.Event{}
	for datetime, id := range i.m.dates {
		if dateCompare(datetime, date) {
			if e, ok := i.m.events[id]; ok {
				events = append(events, eventCreateSafely(e))
			} else {
				return nil, errors.Wrapf(fmt.Errorf("unsynchronize dates and events maps in the inmemory repo"), "can't get event by date: %v", date)
			}
		}
	}
	if len(events) == 0 {
		return nil, domain.ErrEventNotFound
	}
	return events, nil
}

func (i EventRepo) GetForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]*domain.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	events := []*domain.Event{}
	for datetime, id := range i.m.dates {
		if datetime == dateStart || datetime == dateEnd || (datetime.After(dateStart) && datetime.Before(dateEnd)) {
			if e, ok := i.m.events[id]; ok {
				events = append(events, eventCreateSafely(e))
			} else {
				return nil, errors.Wrapf(fmt.Errorf("unsynchronize dates and events maps in the inmemory repo"), "can't get event by period: %v-%v", dateStart, dateEnd)
			}
		}
	}
	if len(events) == 0 {
		return nil, domain.ErrEventNotFound
	}
	return events, nil
}

func (i EventRepo) UpdateByID(ctx context.Context, id string, event domain.Event) error {
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()
	e, ok := i.m.events[id]
	if !ok {
		return domain.ErrEventNotFound
	}
	e.Title = event.Title
	e.DateTime = event.DateTime
	e.UserID = event.UserID
	e.TimeNotify = event.TimeNotify
	e.Duration = event.Duration
	e.Text = event.Text
	return nil
}

func (i EventRepo) DeleteByID(ctx context.Context, id string) error {
	e, err := i.GetByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "can't delete event id:%v", id)
	}
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()

	delete(i.m.events, id)
	delete(i.m.dates, e.DateTime)
	return nil
}

func eventCreateSafely(event *domain.Event) *domain.Event {
	return &domain.Event{
		ID:         event.ID,
		Title:      event.Title,
		DateTime:   event.DateTime,
		Duration:   event.Duration,
		Text:       event.Text,
		UserID:     event.UserID,
		TimeNotify: event.TimeNotify,
	}
}

func dateCompare(dt1, dt2 time.Time) bool {
	y1, m1, d1 := dt1.Date()
	y2, m2, d2 := dt2.Date()
	return (y1 == y2 && m1 == m2 && d1 == d2)
}
