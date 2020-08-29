package app

import (
	"context"
	logger2 "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar_sender/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/queueservice/rabbitservice"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	mainusecase "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

type QueueType string

const (
	rabbit QueueType = "rabbit"
	kafka  QueueType = "kafka"
)

func (a *App) Run(cfg *Config, isDebug bool) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wr := os.Stdout
	if !isDebug {
		wr, err = os.OpenFile(cfg.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrapf(err, "can't create/open log file")
		}
	}

	logger, err := logger2.NewLogger(wr, cfg.Log.Level)
	if err != nil {
		return errors.Wrapf(err, "can't init logger")
	}

	broker, err := InitQueue(cfg, logger)
	if err != nil {
		return errors.Wrapf(err, "can't init queue broker")
	}

	scheduler := usecases.NewSender(broker, logger)

	done := make(chan os.Signal, 2)
	signal.Notify(done, os.Interrupt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info(ctx, "start listen alerts from queue...")
		if err := scheduler.SendingAlerts(ctx); err != nil {
			logger.Error(ctx, errors.Wrapf(err, "can't run sending alerts"))
		}
		logger.Info(ctx, "stop listening alerts")
		done <- os.Interrupt
	}()

	<-done
	cancel()
	wg.Wait()
	return nil
}

func InitQueue(cfg *Config, logger mainusecase.Logger) (entities.NotifyQueue, error) {
	switch cfg.Queue {
	case rabbit:
		broker, err := rabbitservice.NewRabbitManager(cfg.Rabbit.Addr, cfg.Rabbit.ExchangeName, cfg.Rabbit.QueueName, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init db repository")
		}
		return broker, nil
	case kafka:
		return nil, errors.New("kafka broker is not supported")
	default:
		return nil, errors.New("unknown queue type. I know next types:rabbit,kafka")
	}
}
