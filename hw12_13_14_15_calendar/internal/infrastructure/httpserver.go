package infrastructure

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/interfaces"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/usecases"
)

type contextKey string

const RequestID contextKey = "RequestID"

type HTTPServer struct {
	logger  usecases.ILogger
	server  *http.Server
	usecase interfaces.Calendar
	wg      *sync.WaitGroup
}

func NewHTTPServer(wg *sync.WaitGroup, logger usecases.ILogger, usecase interfaces.Calendar) *HTTPServer {
	return &HTTPServer{
		logger:  logger,
		usecase: usecase,
		wg:      wg}
}

type HTTPReqInfo struct {
	ip                    string
	start                 string
	method, path, httpver string
	code                  int
	latency               time.Duration
	useragent             string
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

func (s *HTTPServer) Serve(addr string) {
	defer s.wg.Done()
	handler := interfaces.NewHTTPHandler(s.usecase)

	authMux := http.NewServeMux()
	authMux.HandleFunc("/", handler.Hello)

	authHandler := s.authMiddleware(authMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", authHandler)
	siteMux.HandleFunc("/favicon.ico", handler.Favicon)

	siteHandler := s.accessLogMiddleware(siteMux)
	siteHandler = s.panicMiddleware(siteHandler)

	s.logger.Info(context.Background(), "starting http server at %v", addr)

	s.server = &http.Server{
		Addr:    addr,
		Handler: siteHandler,
	}
	// always returns error. ErrServerClosed on graceful close.
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error(context.Background(), errors.Wrapf(err, "can't start serve server"))
	}
}

func (s *HTTPServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.httpWarn(r.Context(), w, `not authorized user. No userid request value)`, http.StatusUnauthorized)
			return
		}
		userid := r.FormValue("userid")
		if userid == "" {
			s.httpWarn(r.Context(), w, `not authorized user. No userid request value)`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.httpError(r.Context(), w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := setRequestID(r.Context())
		ri := &HTTPReqInfo{
			ip:        ipFromHostPort(r.RemoteAddr),
			start:     start.Format(usecases.LayoutISO),
			method:    r.Method,
			path:      r.URL.Path,
			httpver:   r.Proto,
			useragent: r.Header.Get("User-Agent"),
		}

		rw := &WrapResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r.WithContext(ctx))

		ri.latency = time.Since(start)
		ri.code = rw.status
		s.logger.Info(ctx, "%s [%s] %s %s %s %d %s [%s]", ri.ip, ri.start, ri.method, ri.path, ri.httpver, ri.code, ri.latency, ri.useragent)
	})
}

func (s *HTTPServer) httpError(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Error(ctx, error)
	http.Error(w, error, code)
}

func (s *HTTPServer) httpWarn(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Warn(ctx, error)
	http.Error(w, error, code)
}

type WrapResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *WrapResponseWriter) WriteHeader(status int) {
	w.status = status
}

func setRequestID(ctx context.Context) context.Context {
	reqid := uuid.NewV4()
	return context.WithValue(ctx, RequestID, reqid.String())
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
