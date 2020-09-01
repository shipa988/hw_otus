package httpapi

import (
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"

	"net/http"
)

type Handler interface {
	Index(res http.ResponseWriter, req *http.Request)
	Favicon(res http.ResponseWriter, req *http.Request)
	AddEvent(res http.ResponseWriter, req *http.Request)
	GetDateEvents(res http.ResponseWriter, req *http.Request)
	GetWeekEvents(res http.ResponseWriter, req *http.Request)
	GetMonthEvents(res http.ResponseWriter, req *http.Request)
	UpdateEvent(res http.ResponseWriter, req *http.Request)
	DeleteEvent(res http.ResponseWriter, req *http.Request)
}

type Response interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type ErrorResponse struct {
	Message string
}

type AddResponse struct {
	ID    string         `json:"id"`
	Error *ErrorResponse `json:"error,omitempty"`
}

func NewAddResponse() *AddResponse {
	return &AddResponse{
		ID:    "",
		Error: nil,
	}
}

type IndexResponse struct {
	Message string         `json:"message"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

func NewIndexResponse() *IndexResponse {
	return &IndexResponse{
		Message: "",
		Error:   nil,
	}
}

type GetResponse struct {
	Events []*entities.Event `json:"events"`
	Error  *ErrorResponse    `json:"error,omitempty"`
}

func NewGetResponse() *GetResponse {
	return &GetResponse{
		Events: []*entities.Event{},
		Error:  nil,
	}
}

type UpdateResponse struct {
	ID    string         `json:"id"`
	Error *ErrorResponse `json:"error,omitempty"`
}

func NewUpdateResponse() *UpdateResponse {
	return &UpdateResponse{
		ID:    "",
		Error: nil,
	}
}

type DeleteResponse struct {
	ID    string         `json:"id"`
	Error *ErrorResponse `json:"error,omitempty"`
}

func NewDeleteResponse() *DeleteResponse {
	return &DeleteResponse{
		ID:    "",
		Error: nil,
	}
}
