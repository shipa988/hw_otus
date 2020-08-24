package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/api/httpapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/util"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/mocks"
)

var (
	title    = "title"
	text     = "text"
	dt       = "2020-07-07 12:12:12"
	date     = "2020-07-07"
	duration = "4h"
	tnotify  = "4h"
	userid   = "00112233-4455-6677-8899-aabbccddeeff"
	id       = "00112233-4455-6677-8899-aabbccddeeff"
)

func TestHTTPEndpoints(t *testing.T) {
	testEvent, err := entities.NewEvent(title, dt, duration, text, userid, tnotify)
	testEvent.ID = id
	require.Nil(t, err)
	repo := mocks.NewMockRepo(testEvent)
	logger := mocks.NewMockLogger()
	calendar := usecases.NewCalendar(repo, nil, logger)
	handler := NewAPIHandler(calendar, logger)
	w := httptest.NewRecorder()

	t.Run("AddEvent", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.AddResponse
		}{
			{"good AddEvent",
				httptest.NewRequest(http.MethodPost, "http://127.0.0.1:80/", strings.NewReader(fmt.Sprintf(`title=%v&text=%v&duration=%v&datetime=%v`, title, text, duration, url.QueryEscape(dt)))).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.AddResponse{
					ID:    testEvent.ID,
					Error: nil,
				},
			},
			{"bad title AddEvent",
				httptest.NewRequest(http.MethodGet, "http://127.0.0.1:80/", strings.NewReader(fmt.Sprintf(`title=%v&text=%v&duration=%v&datetime=%v`, title, text, duration, url.QueryEscape(dt)))).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.AddResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "title")},
				},
			},
			{"bad text AddEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?title=%v&duration=%v&datetime=%v", title, duration, url.QueryEscape(dt)), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.AddResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "text")},
				},
			},
			{"bad duration AddEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?text=%v&title=%v&datetime=%v", text, title, url.QueryEscape(dt)), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.AddResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "duration")},
				},
			},
			{"bad datetime AddEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?text=%v&title=%v&duration=%v", text, title, duration), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.AddResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "datetime")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				tcase.r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				handler.AddEvent(w, tcase.r)
				resp := httpapi.NewAddResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.ID, resp.ID)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}

			})

		}
	})
	t.Run("DeleteEvent", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.DeleteResponse
		}{
			{"good DeleteEvent",
				httptest.NewRequest("", fmt.Sprintf("http://127.0.0.1:80/?id=%v", id), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.DeleteResponse{
					ID:    testEvent.ID,
					Error: nil,
				},
			},
			{"bad eventid DeleteEvent",
				httptest.NewRequest("", fmt.Sprintf("http://127.0.0.1:80/"), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.DeleteResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "id")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				handler.DeleteEvent(w, tcase.r)
				resp := httpapi.NewDeleteResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.ID, resp.ID)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}
			})

		}
	})
	t.Run("UpdateEvent", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.UpdateResponse
		}{
			{"good all fields UpdateEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?id=%v&title=%v&text=%v&duration=%v&datetime=%v&timetonotify=%v", id, title, text, duration, url.QueryEscape(dt), tnotify), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.UpdateResponse{
					ID:    testEvent.ID,
					Error: nil,
				},
			},
			{"good at least one field UpdateEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?id=%v&title=%v", id, title), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.UpdateResponse{
					ID:    testEvent.ID,
					Error: nil,
				},
			},
			{"bad no eventid UpdateEvent",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?title=%v", title), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.UpdateResponse{
					ID:    "",
					Error: &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "id")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				handler.UpdateEvent(w, tcase.r)
				resp := httpapi.NewUpdateResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.ID, resp.ID)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}
			})

		}
	})
	t.Run("GetDateEvents", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.GetResponse
		}{
			{"good all fields GetDateEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?date=%v", url.QueryEscape(date)), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{testEvent},
					Error:  nil,
				},
			},
			{"bad no datetime GetDateEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/"), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{},
					Error:  &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "date")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				handler.GetDateEvents(w, tcase.r)
				resp := httpapi.NewGetResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.Events, resp.Events)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}
			})

		}
	})
	t.Run("GetWeekEvents", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.GetResponse
		}{
			{"good all fields GetWeekEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?date=%v", url.QueryEscape(date)), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{testEvent},
					Error:  nil,
				},
			},
			{"bad no datetime GetWeekEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/"), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{},
					Error:  &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "date")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				handler.GetWeekEvents(w, tcase.r)
				resp := httpapi.NewGetResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.Events, resp.Events)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}
			})

		}
	})
	t.Run("GetMonthEvents", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			resp  *httpapi.GetResponse
		}{
			{"good all fields GetMonthEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/?date=%v", url.QueryEscape(date)), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{testEvent},
					Error:  nil,
				},
			},
			{"bad no datetime GetMonthEvents",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:80/"), nil).
					WithContext(util.SetUserID(nil, userid)),
				&httpapi.GetResponse{
					Events: []*entities.Event{},
					Error:  &httpapi.ErrorResponse{Message: fmt.Sprintf(entities.ErrNoField, "date")},
				},
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				handler.GetMonthEvents(w, tcase.r)
				resp := httpapi.NewGetResponse()
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				defer w.Body.Reset()
				require.Nil(t, err)
				require.Equal(t, tcase.resp.Events, resp.Events)
				if resp.Error != nil && tcase.resp.Error == nil {
					t.Errorf("error must be nil ")
				}
				if resp.Error == nil && tcase.resp.Error != nil {
					t.Errorf("error must be not nil ")
				}
			})

		}
	})
}
