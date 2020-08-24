package app

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar_scheduler/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/queueservice/rabbitservice"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/repository/db"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/repository/inmemory"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	mainusecase "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

type RepoType string

const (
	dbRepo       RepoType = "db"
	inmemoryRepo RepoType = "inmemory"
)

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

	logger, err := internal.NewLogger(wr, cfg.Log.Level)
	if err != nil {
		return errors.Wrapf(err, "can't init logger")
	}

	repo, err := InitRepo(cfg, logger)
	if err != nil {
		return errors.Wrapf(err, "can't init repository")
	}

	broker, err := InitQueue(cfg, logger)
	if err != nil {
		return errors.Wrapf(err, "can't init queue broker")
	}

	scheduler := usecases.NewScheduler(repo, broker, logger)

	done := make(chan os.Signal, 3)
	signal.Notify(done, os.Interrupt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info(ctx, "start sending alerts to queue...")
		scheduler.SendAlerts(ctx)
		logger.Info(ctx, "stop sending alerts to queue...")
		done <- os.Interrupt
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info(ctx, "start cleaning old events from repo...")
		scheduler.CleanEvents(ctx, 360)
		logger.Info(ctx, "stop cleaning old events from repo...")
		done <- os.Interrupt
	}()
	<-done
	cancel()
	wg.Wait()
	return nil
}

func InitRepo(cfg *Config, logger mainusecase.Logger) (entities.EventRepo, error) {
	switch cfg.RepoType {
	case dbRepo:
		repo, err := db.NewDBEventRepo(cfg.DB.Driver, cfg.DB.DSN, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init db repository")
		}
		return repo, nil
	case inmemoryRepo:
		repo, err := inmemory.NewInMemoryEventRepo(inmemory.NewMapRepo(), logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init inmemo repository")
		}
		return repo, nil
	default:
		return nil, errors.New("unknown repository type. I know next types:db-database,inmemory-map struct into app")
	}
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
