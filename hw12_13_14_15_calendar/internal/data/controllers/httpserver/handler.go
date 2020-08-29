package httpserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	api "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/api/httpapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/util"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

// HTTPHandler preparing http.Handler with all middlewares. I extract this method for test whole the server. What another way can I test whole server?
type HTTPHandler struct {
	logger     usecases.Logger
	apiHandler api.Handler
}

func GetHandler(logger usecases.Logger, apiHandler api.Handler) http.Handler {
	h := &HTTPHandler{
		logger:     logger,
		apiHandler: apiHandler,
	}
	authMux := http.NewServeMux()
	authMux.HandleFunc("/", h.get(h.apiHandler.Index))
	authMux.HandleFunc("/events/add", h.post(h.apiHandler.AddEvent))
	authMux.HandleFunc("/events/date", h.get(h.apiHandler.GetDateEvents))
	authMux.HandleFunc("/events/week", h.get(h.apiHandler.GetWeekEvents))
	authMux.HandleFunc("/events/month", h.get(h.apiHandler.GetMonthEvents))
	authMux.HandleFunc("/events/update", h.patch(h.apiHandler.UpdateEvent))
	authMux.HandleFunc("/events/delete", h.delete(h.apiHandler.DeleteEvent))

	authHandler := h.authMiddleware(authMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", authHandler)
	siteMux.HandleFunc("/favicon.ico", h.apiHandler.Favicon)

	siteHandler := h.accessLogMiddleware(siteMux)
	siteHandler = h.panicMiddleware(siteHandler)
	return siteHandler
}

func (s *HTTPHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid := r.Header.Get("x-user-id")
		if len(userid) == 0 {
			s.httpWarn(r.Context(), w, `not authorized user`, http.StatusUnauthorized)
			return
		}
		ctx := util.SetUserID(r.Context(), userid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *HTTPHandler) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.httpError(r.Context(), w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPHandler) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := util.SetRequestID(r.Context())
		ri := util.NewHTTPReqInfo(r.RemoteAddr, start, r.Method, r.URL.Path, r.Proto, r.Header.Get("User-Agent"))
		rw := &WrapResponseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r.WithContext(ctx))

		ri.Latency = time.Since(start)
		ri.Code = strconv.Itoa(rw.status)
		s.logRequest(ctx, ri)
	})
}

func (s *HTTPHandler) get(handler http.HandlerFunc) http.HandlerFunc {
	return s.methods(handler, http.MethodGet)
}

func (s *HTTPHandler) post(handler http.HandlerFunc) http.HandlerFunc {
	return s.methods(handler, http.MethodPost)
}

func (s *HTTPHandler) delete(handler http.HandlerFunc) http.HandlerFunc {
	return s.methods(handler, http.MethodDelete)
}

func (s *HTTPHandler) patch(handler http.HandlerFunc) http.HandlerFunc {
	return s.methods(handler, http.MethodPatch)
}

func (s *HTTPHandler) methods(handler http.HandlerFunc, methodName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == methodName {
			handler.ServeHTTP(w, req)
			return
		}
		s.httpWarn(req.Context(), w, `not supported method: `+req.Method, http.StatusMethodNotAllowed)
	}
}

func (s *HTTPHandler) logRequest(ctx context.Context, ri *util.HTTPReqInfo) {
	s.logger.Info(ctx, "%s [%s] %s %s %s %s %s [%s]", ri.IP, ri.Start, ri.Method, ri.Path, ri.Httpver, ri.Code, ri.Latency, ri.Useragent)
}

func (s *HTTPHandler) httpError(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Error(ctx, error)
	http.Error(w, error, code)
}

func (s *HTTPHandler) httpWarn(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Warn(ctx, error)
	http.Error(w, error, code)
}

type WrapResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *WrapResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}
