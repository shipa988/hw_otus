package usecases

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

const (
	ErrSend  = "can't send alert"
	ErrClean = "can't clean events"
)

var _ Scheduler = (*SchedulerInteractor)(nil)

type SchedulerInteractor struct {
	events    entities.EventRepo
	alerts    entities.NotifyQueue
	logger    usecases.Logger
	alertCash alertCash
}

func NewScheduler(eRepo entities.EventRepo, nQueue entities.NotifyQueue, logger usecases.Logger) *SchedulerInteractor {
	return &SchedulerInteractor{
		events: eRepo,
		alerts: nQueue,
		logger: logger,
		alertCash: alertCash{
			cash:    map[string]time.Time{},
			RWMutex: &sync.RWMutex{},
		},
	}
}

func (s *SchedulerInteractor) SendAlerts(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.alertCash.cleanCash(ctx)
		s.logger.Info(ctx, "stop cleaning alerts from cash")
	}()
	loop := true
	tick := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		loop = false
	default:
		s.doSend(ctx)
	}
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			tick.Stop()
			break
		case <-tick.C: //todo: more correct is notify about all events not be alerting before(eg. when scheduler is off), and add flag "isnotified" in DB, but i decided add local day-alerts alertCash
			s.doSend(ctx)
		}
	}
	wg.Wait()
}

func (s *SchedulerInteractor) doSend(ctx context.Context) {
	events, e := s.events.GetByNotifyDate(ctx, time.Now())
	if e != nil {
		s.logger.Error(ctx, errors.Wrap(e, ErrSend))
	}
	for _, event := range events {
		s.alertCash.RLock()
		if _, ok := s.alertCash.cash[event.ID]; !ok {
			s.alertCash.cash[event.ID] = event.DateTime
			s.logger.Info(context.TODO(), fmt.Sprintf(`new event "id:%v title:%v" for alerting is found`, event.ID, event.Title))
			n := entities.Notify{
				ID:       event.ID,
				Title:    event.Title,
				UserID:   event.UserID,
				DateTime: event.DateTime,
			}

			go func() {
				if err := s.alerts.Push(n); err != nil {
					s.logger.Error(context.TODO(), errors.Wrapf(err, ErrSend))
				}
			}()
		}
		s.alertCash.RUnlock()
	}
}

func (s *SchedulerInteractor) CleanEvents(ctx context.Context, daysToDelete int) {
	loop := true
	tick := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		loop = false
	default:
		s.doClean(ctx, daysToDelete)
	}
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			tick.Stop()
			break
		case <-tick.C:
			s.doClean(ctx, daysToDelete)
		}
	}
}

func (s *SchedulerInteractor) doClean(ctx context.Context, daysToDelete int) {
	events, err := s.events.GetForPeriod(ctx, time.Time{}, time.Now().AddDate(0, 0, -daysToDelete))
	if err != nil {
		s.logger.Error(context.TODO(), errors.Wrapf(err, ErrClean))
	}
	for _, event := range events {
		err := s.events.DeleteByID(ctx, event.ID)
		if err != nil {
			s.logger.Error(context.TODO(), errors.Wrapf(err, ErrClean))
		}
		s.logger.Info(context.TODO(), "old event %v is deleted", event.ID)
	}
}

type alertCash struct {
	cash map[string]time.Time
	*sync.RWMutex
}

func (a *alertCash) cleanCash(ctx context.Context) {
	loop := true
	tick := time.NewTicker(time.Minute)
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			tick.Stop()
			break
		case <-tick.C:
			a.Lock()
			for id, ntime := range a.cash {
				if ntime.Before(time.Now().AddDate(0, 0, -1)) {
					delete(a.cash, id)
				}
			}
			a.Unlock()
		}
	}
}
