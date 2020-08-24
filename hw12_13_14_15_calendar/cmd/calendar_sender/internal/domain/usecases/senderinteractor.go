package usecases

import (
	"context"
	"fmt"
	"sync"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

const (
	ErrSend  = "can't send alert"
	ErrClean = "can't clean events"
)

var _ Sender = (*SenderInteractor)(nil)

type SenderInteractor struct {
	alerts entities.NotifyQueue
	logger usecases.Logger
}

func NewSender(nQueue entities.NotifyQueue, logger usecases.Logger) *SenderInteractor {
	return &SenderInteractor{
		alerts: nQueue,
		logger: logger,
	}
}

func (s *SenderInteractor) processAlert(ctx context.Context, alerts <-chan entities.Notify) {
	loop := true
	for loop {
		select {
		case alert, ok := <-alerts:
			if !ok {
				loop = false
				break
			}
			fmt.Printf("%+v\n", alert)
		case <-ctx.Done():
			loop = false
			break
		}
	}
}
func (s *SenderInteractor) SendingAlerts(ctx context.Context) error {
	notifies := make(chan entities.Notify)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processAlert(ctx, notifies)
	}()
	err := s.alerts.Pull(ctx, notifies)
	if err != nil {
		close(notifies)
	}
	wg.Wait()
	return err
}
