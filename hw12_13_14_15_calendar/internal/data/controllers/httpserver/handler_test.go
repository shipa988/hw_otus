package httpserver

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/mocks"
)

const (
	ErrUnAuthorize        = "not authorized user\n"
	ErrNotSupportedMethod = "not supported method: "
)

func TestHTTPHandler(t *testing.T) {
	apiHandler := mocks.NewMockHandler()
	logger := mocks.NewMockLogger()
	httpHandler := GetHandler(logger, apiHandler)
	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	client := &http.Client{}
	testCases := []struct {
		name    string
		method  string
		url     string
		auth    bool
		expResp string
		expCode int
		err     error
	}{
		{"add: good ", http.MethodPost, "events/add", true, "AddEvent", 200, nil},
		{"add: unauthorized", http.MethodPost, "events/add", false, ErrUnAuthorize, 401, nil},
		{"add: not supported method", http.MethodGet, "events/add", true, ErrNotSupportedMethod + "GET\n", 405, nil},

		{"delete: good ", http.MethodDelete, "events/delete", true, "DeleteEvent", 200, nil},
		{"delete: unauthorized", http.MethodPost, "events/delete", false, ErrUnAuthorize, 401, nil},
		{"delete: not supported method", http.MethodGet, "events/delete", true, ErrNotSupportedMethod + "GET\n", 405, nil},

		{"update: good ", http.MethodPatch, "events/update", true, "UpdateEvent", 200, nil},
		{"update: unauthorized", http.MethodPost, "events/update", false, ErrUnAuthorize, 401, nil},
		{"update: not supported method", http.MethodGet, "events/update", true, ErrNotSupportedMethod + "GET\n", 405, nil},

		{"getdate: good ", http.MethodGet, "events/date", true, "GetDateEvents", 200, nil},
		{"getdate: unauthorized", http.MethodGet, "events/date", false, ErrUnAuthorize, 401, nil},
		{"getdate: not supported method", http.MethodPost, "events/date", true, ErrNotSupportedMethod + "POST\n", 405, nil},

		{"getweek: good ", http.MethodGet, "events/week", true, "GetWeekEvents", 200, nil},
		{"getweek: unauthorized", http.MethodGet, "events/week", false, ErrUnAuthorize, 401, nil},
		{"getweek: not supported method", http.MethodPost, "events/week", true, ErrNotSupportedMethod + "POST\n", 405, nil},

		{"getmonth: good ", http.MethodGet, "events/month", true, "GetMonthEvents", 200, nil},
		{"getmonth: unauthorized", http.MethodGet, "events/month", false, ErrUnAuthorize, 401, nil},
		{"getmonth: not supported method", http.MethodPost, "events/month", true, ErrNotSupportedMethod + "POST\n", 405, nil},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			req, err := http.NewRequest(tCase.method, strings.Join([]string{testServer.URL, tCase.url}, "/"), nil)
			if tCase.auth {
				req.Header.Add("x-user-id", "123123123")
			}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			require.Equal(t, tCase.err, err)
			b := bytes.Buffer{}
			io.Copy(&b, resp.Body)
			require.Equal(t, tCase.expResp, b.String())
			require.Equal(t, tCase.expCode, resp.StatusCode)
		})
	}

}
