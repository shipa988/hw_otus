//nolint: errcheck
package mocks

import (
	"net/http"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/api/httpapi"
)

var _ httpapi.Handler = (*MockHandler)(nil)

type MockHandler struct{}

func (m MockHandler) Index(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Index"))
}

func (m MockHandler) Favicon(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Favicon"))
}

func (m MockHandler) AddEvent(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("AddEvent"))
}

func (m MockHandler) GetDateEvents(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("GetDateEvents"))
}

func (m MockHandler) GetWeekEvents(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("GetWeekEvents"))
}

func (m MockHandler) GetMonthEvents(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("GetMonthEvents"))
}

func (m MockHandler) UpdateEvent(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("UpdateEvent"))
}

func (m MockHandler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("DeleteEvent"))
}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}
