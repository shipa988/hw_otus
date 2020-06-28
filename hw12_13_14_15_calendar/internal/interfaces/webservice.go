package interfaces

import (
	"fmt"
	"net/http"
)

type HTTPHandler struct {
	calendar Calendar
}

func NewHTTPHandler(calendar Calendar) *HTTPHandler {
	return &HTTPHandler{
		calendar: calendar}
}

//nolint:errcheck
func (handler HTTPHandler) Hello(res http.ResponseWriter, req *http.Request) {
	_ = req.Context() //todo: don't forget use ctx for usecases (handler.calendar.DeleteEvent(ctx,id))

	err := req.ParseForm()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userid := req.FormValue("userid")
	res.Write([]byte(fmt.Sprintf("hello %v!", userid)))
	res.WriteHeader(http.StatusOK)
}

func (handler HTTPHandler) Favicon(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "image/x-icon")
	res.Header().Set("Cache-Control", "public, max-age=7776000")
	res.WriteHeader(http.StatusOK)
}
