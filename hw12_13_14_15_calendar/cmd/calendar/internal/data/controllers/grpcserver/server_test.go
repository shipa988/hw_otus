package grpcserver

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	api "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/api/grpcapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/mocks"
)

const buffer = 1024 * 1024

var (
	title   = "title"
	text    = "text"
	dt      = "2020-07-07 12:12:12"
	date    = "2020-07-07"
	dur     = "4h"
	tnotify = "4h"
	userid  = "00112233-4455-6677-8899-aabbccddeeff"
	id      = "00112233-4455-6677-8899-aabbccddeeff"
)

func TestGRPCServer(t *testing.T) {
	// init test.
	wg := &sync.WaitGroup{}
	testEvent, err := entities.NewEvent(title, dt, dur, text, userid, tnotify)
	testEvent.ID = id
	require.Nil(t, err)

	ptypesTs, err := types.TimestampProto(testEvent.DateTime)
	require.Nil(t, err)
	pbTs := &timestamp.Timestamp{
		Seconds: ptypesTs.Seconds,
		Nanos:   ptypesTs.Nanos,
	}

	ptypesDu := types.DurationProto(testEvent.Duration)
	pbDu := &duration.Duration{
		Seconds: ptypesDu.Seconds,
		Nanos:   ptypesDu.Nanos,
	}

	ptypesTn := types.DurationProto(testEvent.TimeNotify)
	pbTn := &duration.Duration{
		Seconds: ptypesTn.Seconds,
		Nanos:   ptypesTn.Nanos,
	}

	testPBEvent := &api.Event{
		Id:         testEvent.ID,
		Title:      testEvent.Title,
		Datetime:   pbTs,
		Duration:   pbDu,
		Text:       testEvent.Text,
		Userid:     testEvent.UserID,
		Timenotify: pbTn,
	}
	testPBEvents := &api.Events{Event: []*api.Event{testPBEvent}}

	repo := mocks.NewMockRepo(testEvent)
	logger := mocks.NewMockLogger()
	calendar := usecases.NewCalendar(repo, nil, logger)

	server := NewGRPCServer(wg, logger, calendar)
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	defer func() {
		listener.Close()
		server.StopServe()
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(ctx context.Context, s string) (conn net.Conn, err error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	client := api.NewCalendarServiceClient(conn)

	wg.Add(1)
	go server.Serve(listener)

	// run tests.
	t.Run("add: ok", func(t *testing.T) {
		req := &api.AddEventRequest{Title: title, Text: text, Datetime: dt, Duration: dur}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.AddEvent(ctx, req)

		expresp := &api.AddEventResponse{
			Result: &api.AddEventResponse_Id{Id: id},
		}

		require.Nil(t, err)
		require.Equal(t, expresp.GetId(), r.GetId())
	})
	t.Run("add: unauthorized add", func(t *testing.T) {
		req := &api.AddEventRequest{Title: title, Text: text, Datetime: dt, Duration: dur}

		r, err := client.AddEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("add: bad Datetime field", func(t *testing.T) {
		req := &api.AddEventRequest{Title: title, Text: text, Datetime: "bad", Duration: dur}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.AddEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	t.Run("delete: ok", func(t *testing.T) {
		req := &api.DeleteEventRequest{Id: id}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.DeleteEvent(ctx, req)

		expresp := &api.DeleteEventResponse{
			Result: &api.DeleteEventResponse_Id{Id: id},
		}

		require.Nil(t, err)
		require.Equal(t, expresp.GetId(), r.GetId())
	})
	t.Run("delete: unauthorized add", func(t *testing.T) {
		req := &api.DeleteEventRequest{Id: id}

		r, err := client.DeleteEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("delete: bad no id field", func(t *testing.T) {
		req := &api.DeleteEventRequest{}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.DeleteEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	t.Run("update: ok", func(t *testing.T) {
		req := &api.UpdateEventRequest{Id: id, Title: "new title"}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.UpdateEvent(ctx, req)

		expresp := &api.UpdateEventResponse{
			Result: &api.UpdateEventResponse_Id{Id: id},
		}

		require.Nil(t, err)
		require.Equal(t, expresp.GetId(), r.GetId())
	})
	t.Run("update: unauthorized add", func(t *testing.T) {
		req := &api.UpdateEventRequest{Id: id, Title: title, Text: text, Datetime: dt, Duration: dur}

		r, err := client.UpdateEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("update: bad no id field", func(t *testing.T) {
		req := &api.UpdateEventRequest{Title: title}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", "1"))

		r, err := client.UpdateEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	t.Run("getdate: ok", func(t *testing.T) {
		req := &api.GetDateEventRequest{Date: date}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.GetDateEvent(ctx, req)

		expresp := &api.GetDateEventResponse{
			Result: &api.GetDateEventResponse_Events{testPBEvents},
		}

		require.Nil(t, err)
		require.EqualValues(t, expresp.GetEvents(), r.GetEvents())
	})
	t.Run("getdate: unauthorized add", func(t *testing.T) {
		req := &api.GetDateEventRequest{Date: date}

		r, err := client.GetDateEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("getdate: bad no Date field", func(t *testing.T) {
		req := &api.GetDateEventRequest{}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", "1"))

		r, err := client.GetDateEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	t.Run("getweek: ok", func(t *testing.T) {
		req := &api.GetWeekEventRequest{Date: date}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.GetWeekEvent(ctx, req)

		expresp := &api.GetWeekEventResponse{
			Result: &api.GetWeekEventResponse_Events{testPBEvents},
		}

		require.Nil(t, err)
		require.EqualValues(t, expresp.GetEvents(), r.GetEvents())
	})
	t.Run("getweek: unauthorized add", func(t *testing.T) {
		req := &api.GetWeekEventRequest{Date: date}

		r, err := client.GetWeekEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("getweek: bad no Date field", func(t *testing.T) {
		req := &api.GetWeekEventRequest{}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", "1"))

		r, err := client.GetWeekEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	t.Run("getmonth: ok", func(t *testing.T) {
		req := &api.GetMonthEventRequest{Date: date}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", userid))

		r, err := client.GetMonthEvent(ctx, req)

		expresp := &api.GetMonthEventResponse{
			Result: &api.GetMonthEventResponse_Events{testPBEvents},
		}

		require.Nil(t, err)
		require.EqualValues(t, expresp.GetEvents(), r.GetEvents())
	})
	t.Run("getmonth: unauthorized add", func(t *testing.T) {
		req := &api.GetMonthEventRequest{Date: date}

		r, err := client.GetMonthEvent(context.Background(), req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.PermissionDenied.String())
	})
	t.Run("getmonth: bad no Date field", func(t *testing.T) {
		req := &api.GetMonthEventRequest{}
		//auth
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-user-id", "1"))

		r, err := client.GetMonthEvent(ctx, req)

		require.Nil(t, r)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), codes.Aborted.String())
	})

	server.StopServe()
	wg.Wait()
}
