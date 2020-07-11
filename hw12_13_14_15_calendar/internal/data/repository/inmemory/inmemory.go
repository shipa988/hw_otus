package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

var _ entities.EventRepo = (*EventRepo)(nil)

type dates map[time.Time]*entities.Event
type MapRepo struct {
	rwmux  *sync.RWMutex
	users  map[string]dates
	events map[string]*entities.Event
}

func NewMapRepo() *MapRepo {
	return &MapRepo{
		rwmux:  &sync.RWMutex{},
		users:  make(map[string]dates),
		events: make(map[string]*entities.Event),
	}
}
func (m *MapRepo) Clear() {
	m.rwmux.Lock()
	defer m.rwmux.Unlock()
	m.users = make(map[string]dates)
	m.events = make(map[string]*entities.Event)
}

type EventRepo struct {
	m      *MapRepo //emulate infrastructure layer over type MapRepo
	logger usecases.Logger
}

func NewInMemoryEventRepo(m *MapRepo, logger usecases.Logger) (*EventRepo, error) {
	return &EventRepo{
		m:      m,
		logger: logger,
	}, nil
}

func (i EventRepo) Add(ctx context.Context, event entities.Event) (string, error) {
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()
	if dates, ok := i.m.users[event.UserID]; !ok {
		//user not exist
		i.m.users[event.UserID] = make(map[time.Time]*entities.Event)
	} else if _, ok := dates[event.DateTime]; ok {
		//date is busy
		return "", entities.ErrDateBusy
	}
	id := uuid.NewV4().String()
	event.ID = id
	i.m.users[event.UserID][event.DateTime] = &event
	i.m.events[event.ID] = &event
	return id, nil
}

func (i EventRepo) GetByID(ctx context.Context, userID, eventID string) (*entities.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	_, ok := i.m.users[userID]
	if !ok {
		return nil, errors.Wrapf(entities.ErrUnknownUser, "can't get event by id: %v", eventID)
	}
	if event, ok := i.m.events[eventID]; ok {
		return eventCreateSafely(event), nil
	}
	return nil, entities.ErrEventNotFound
}

func (i EventRepo) GetByDate(ctx context.Context, userID string, date time.Time) ([]*entities.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	events := []*entities.Event{}
	dates, ok := i.m.users[userID]
	if !ok {
		return nil, errors.Wrapf(entities.ErrUnknownUser, "can't get event by date: %v", date)
	}
	for datetime, event := range dates {
		if dateCompare(datetime, date) {
			events = append(events, eventCreateSafely(event))
		}
	}
	if len(events) == 0 {
		return nil, entities.ErrEventNotFound
	}
	return events, nil
}

func (i EventRepo) GetForPeriod(ctx context.Context, userID string, dateStart time.Time, dateEnd time.Time) ([]*entities.Event, error) {
	i.m.rwmux.RLock()
	defer i.m.rwmux.RUnlock()
	events := []*entities.Event{}
	dates, ok := i.m.users[userID]
	if !ok {
		return nil, errors.Wrapf(entities.ErrUnknownUser, "can't get event by period: %v-%v", dateStart, dateEnd)
	}
	for datetime, event := range dates {
		if datetime == dateStart || datetime == dateEnd || (datetime.After(dateStart) && datetime.Before(dateEnd)) {
			events = append(events, eventCreateSafely(event))
		}
	}
	if len(events) == 0 {
		return nil, entities.ErrEventNotFound
	}
	return events, nil
}

func (i EventRepo) UpdateByID(ctx context.Context, userID, eventID string, event entities.Event) error {
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()
	_, ok := i.m.users[userID]
	if !ok {
		return errors.Wrapf(entities.ErrUnknownUser, "can't update event by id %v and userid %v", eventID, userID)
	}
	e, ok := i.m.events[eventID]
	if !ok {
		return entities.ErrEventNotFound
	}
	e.Title = event.Title
	e.DateTime = event.DateTime
	e.UserID = event.UserID
	e.TimeNotify = event.TimeNotify
	e.Duration = event.Duration
	e.Text = event.Text
	return nil
}

func (i EventRepo) DeleteByID(ctx context.Context, userID, eventID string) error {
	e, err := i.GetByID(ctx, userID, eventID)
	if err != nil {
		return errors.Wrapf(err, "can't delete event id:%v", eventID)
	}
	i.m.rwmux.Lock()
	defer i.m.rwmux.Unlock()

	delete(i.m.events, eventID)
	delete(i.m.users[userID], e.DateTime)
	return nil
}

func eventCreateSafely(event *entities.Event) *entities.Event {
	return &entities.Event{
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
	return dt1 == dt2
	/*y1, m1, d1 := dt1.Date()
	y2, m2, d2 := dt2.Date()
	return (y1 == y2 && m1 == m2 && d1 == d2)*/
}
