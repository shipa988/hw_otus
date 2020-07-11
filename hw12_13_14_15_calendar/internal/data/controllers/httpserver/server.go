package httpserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

type HTTPServer struct {
	logger      usecases.Logger
	server      *http.Server
	httpHandler http.Handler
	wg          *sync.WaitGroup
}

func NewHTTPServer(wg *sync.WaitGroup, logger usecases.Logger, httpHandler http.Handler) *HTTPServer {
	return &HTTPServer{
		logger:      logger,
		httpHandler: httpHandler,
		wg:          wg}
}

func (s *HTTPServer) Serve(addr string) {
	defer s.wg.Done()

	s.logger.Info(context.Background(), "starting http server at %v", addr)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.httpHandler,
	}

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error(context.Background(), errors.Wrapf(err, "can't start http server at %v", addr))
	}
}

func (s *HTTPServer) StopServe() {
	ctx := context.Background()
	s.logger.Info(ctx, "stopping http server")
	defer s.logger.Info(ctx, "http server stopped")
	if s.server == nil {
		s.logger.Error(ctx, "http server is nil")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(ctx, "can't stop http server with error: %v", err)
	}
}
