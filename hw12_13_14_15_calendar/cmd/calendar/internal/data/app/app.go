package app

import (
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/data/controllers/grpcserver"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/data/controllers/httpserver"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/repository/db"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/repository/inmemory"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	mainusecase "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/logger"

	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(cfg *Config, isDebug bool) (err error) {
	wr := os.Stdout
	if !isDebug {
		wr, err = os.OpenFile(cfg.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrapf(err, "can't create/open log file")
		}
	}

	logger, err := logger.NewLogger(wr, cfg.Log.Level)
	if err != nil {
		return errors.Wrapf(err, "can't init logger")
	}

	repo, err := InitRepo(cfg, logger)
	if err != nil {
		return errors.Wrapf(err, "can't init repository")
	}
	calendar := usecases.NewCalendar(repo, nil, logger)
	// set executors for api.
	apiHandler := httpserver.NewAPIHandler(calendar, logger)
	// prepare http handler with all middlewares for server.
	httpHandler := httpserver.GetHandler(logger, apiHandler)

	wg := &sync.WaitGroup{}
	// prepare http server with handler
	httpServer := httpserver.NewHTTPServer(wg, logger, httpHandler)
	// prepare grpc server with handler
	grpcServer := grpcserver.NewGRPCServer(wg, logger, calendar)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	wg.Add(1)
	go httpServer.Serve(net.JoinHostPort("0.0.0.0", cfg.API.HTTPPort))

	wg.Add(1)
	l := grpcServer.PrepareGRPCListener(net.JoinHostPort("0.0.0.0", cfg.API.GRPCPort))
	go grpcServer.Serve(l)

	wg.Add(1)
	go grpcServer.ServeGW(net.JoinHostPort("0.0.0.0", cfg.API.GRPCPort), net.JoinHostPort("0.0.0.0", cfg.API.GRPCGWPort))

	go func() {
		wg.Wait()
		quit <- os.Interrupt
	}()

	<-quit

	httpServer.StopServe()
	grpcServer.StopServe()
	grpcServer.StopGWServe()

	return nil
}

func InitRepo(cfg *Config, logger mainusecase.Logger) (entities.EventRepo, error) {
	switch cfg.RepoType {
	case "db":
		repo, err := db.NewDBEventRepo(cfg.DB.Driver, cfg.DB.DSN, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init db repository")
		}
		return repo, nil
	case "inmemory":
		repo, err := inmemory.NewInMemoryEventRepo(inmemory.NewMapRepo(), logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't init inmemo repository")
		}
		return repo, nil
	default:
		return nil, errors.New("unknown repository type. I know next types:db-database,inmemory-map struct into app")
	}
}
