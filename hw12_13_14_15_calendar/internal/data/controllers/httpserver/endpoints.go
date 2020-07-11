//nolint:errcheck,unparam
package httpserver

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/api/httpapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/util"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

const (
	ErrParseParams = "can't parse request params"
)

var _ api.Handler = (*APIHandler)(nil)

type APIHandler struct {
	calendar usecases.Calendar
	logger   usecases.Logger
}

func NewAPIHandler(calendar usecases.Calendar, logger usecases.Logger) *APIHandler {
	return &APIHandler{
		calendar: calendar,
		logger:   logger}
}

func (handler APIHandler) Index(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewIndexResponse()
	code := http.StatusOK

	defer func() {
		msg, err := resp.MarshalJSON()
		if err != nil {
			code = http.StatusInternalServerError
		}
		res.Write(msg)
		res.WriteHeader(code)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}

	resp.Message = fmt.Sprintf("hello %v!", userid)
}

func (handler APIHandler) Favicon(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "image/x-icon")
	res.Header().Set("Cache-Control", "public, max-age=7776000")
	res.WriteHeader(http.StatusOK)
}

func (handler APIHandler) AddEvent(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewAddResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}
	title := handler.getParam(req, "title")
	text := handler.getParam(req, "text")
	datetime := handler.getParam(req, "datetime")
	duration := handler.getParam(req, "duration")
	timetonotify := handler.getParam(req, "timetonotify")

	id, err := handler.calendar.MakeEvent(ctx, title, datetime, text, userid, duration, timetonotify)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.ID = id
}

func (handler APIHandler) GetDateEvents(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewGetResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}

	date := handler.getParam(req, "date")

	events, err := handler.calendar.GetDateEvents(ctx, date, userid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.Events = events
}

func (handler APIHandler) GetWeekEvents(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewGetResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}
	date := handler.getParam(req, "date")

	events, err := handler.calendar.GetWeekEvents(ctx, date, userid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.Events = events
}

func (handler APIHandler) GetMonthEvents(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewGetResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}

	date := handler.getParam(req, "date")

	events, err := handler.calendar.GetMonthEvents(ctx, date, userid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.Events = events
}

func (handler APIHandler) UpdateEvent(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewUpdateResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}

	eventid := handler.getParam(req, "id")
	title := handler.getParam(req, "title")
	text := handler.getParam(req, "text")
	datetime := handler.getParam(req, "datetime")
	duration := handler.getParam(req, "duration")
	timetonotify := handler.getParam(req, "timetonotify")

	id, err := handler.calendar.UpdateEvent(ctx, userid, eventid, title, datetime, text, duration, timetonotify)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.ID = id
}

func (handler APIHandler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userid := util.GetUserID(ctx)

	resp := api.NewDeleteResponse()
	code := http.StatusOK
	// defer marshal answer and return response with code.
	defer func() {
		handler.sendResponse(resp, code, res)
	}()
	// parse all request params.
	err := req.ParseForm()
	if err != nil {
		code = http.StatusBadRequest
		handler.info(ctx, ErrParseParams)
		return
	}

	eventid := handler.getParam(req, "id")

	id, err := handler.calendar.DeleteEvent(ctx, userid, eventid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Error = handler.error(ctx, err)
		return
	}
	resp.ID = id
}

func (handler APIHandler) error(ctx context.Context, error error) *api.ErrorResponse {
	handler.logger.Error(ctx, error)
	return &api.ErrorResponse{Message: fmt.Sprint(error.Error())}
}

func (handler APIHandler) info(ctx context.Context, info string, args ...interface{}) *api.ErrorResponse {
	handler.logger.Info(ctx, info, args...)
	return &api.ErrorResponse{Message: fmt.Sprintf(info, args...)}
}

func (handler APIHandler) getParam(req *http.Request, name string) string {
	params, ok := req.Form[name]
	if !ok || len(params) == 0 || len(params[0]) < 1 {
		return ""
	}
	return params[0]
}

func (handler APIHandler) sendResponse(resp api.Response, code int, res http.ResponseWriter) {
	msg, err := resp.MarshalJSON()
	if err != nil {
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(msg)
}
