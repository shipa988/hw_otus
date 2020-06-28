package infrastructure

import (
	"github.com/pkg/errors"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/interfaces/repository/db"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/interfaces/repository/inmemory"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/usecases"

	"net"
	"os"
	"os/signal"
	"sync"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(cfg *AppConfig, isDebug bool) (err error) {
	wr := os.Stdout
	if !isDebug {
		wr, err = os.OpenFile(cfg.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrapf(err, "can't create/open log file")
		}
	}

	logger, err := NewLogger(wr, cfg.Log.Level)
	if err != nil {
		return errors.Wrapf(err, "can't init logger")
	}

	repo, err := InitRepo(cfg, logger)
	if err != nil {
		return errors.Wrapf(err, "can't init repository")
	}
	calendar := usecases.NewCalendar(repo, nil, logger)

	wg := &sync.WaitGroup{}
	server := NewHTTPServer(wg, logger, calendar)
	wg.Add(1)
	go server.Serve(net.JoinHostPort("localhost", cfg.API.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.StopServe()
	return nil
}

func InitRepo(cfg *AppConfig, logger usecases.ILogger) (domain.EventRepo, error) {
	switch cfg.RepoType {
	case "db":
		repo, err := db.NewDBEventRepo(cfg.DB.Driver, cfg.DB.DSN, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init db repository")
		}
		return repo, nil
	case "inmemo":
		repo, err := inmemory.NewInMemoryEventRepo(inmemory.NewMapRepo(), logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init inmemo repository")
		}
		return repo, nil
	default:
		return nil, errors.New("unknown repository type. I know next types:db-database,inmemory-map struct into app")
	}
}
