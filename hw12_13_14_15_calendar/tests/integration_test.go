// +build integration

package tests

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	grpcserver "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/api/grpcapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"log"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/repository/db"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/logger"
)

var (
	title   = "title"
	text    = "text"
	dt      = "2020-07-07 12:12:12"
	date    = "2020-07-07"
	dur     = "4h"
	tnotify = "24h"
	userid  = "75e45c9c-5365-45fc-b6f9-8343c438666c"
	id      = "00112233-4455-6677-8899-aabbccddeeff"
)

var (
	validMetadata   = metadata.New(map[string]string{"x-user-id": "75e45c9c-5365-45fc-b6f9-8343c438666c"})
	invalidMetadata = metadata.New(map[string]string{"x-user-id-invalid": "11111"})
)

type Suite struct {
	suite.Suite
	client    grpcserver.CalendarServiceClient
	testEvent *entities.Event
	conn      *grpc.ClientConn
	repo      *db.EventRepo
}
type condition func()

func TestIntegration(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

func (s *Suite) SetupSuite() {
	//port := os.Getenv("GRPC_PORT")
	//dsn := os.Getenv("DSN")
	port := "4445"
	dsn := "host=localhost port=5432 user=igor password=igor dbname=calendar sslmode=disable"
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	serverAddr := net.JoinHostPort("localhost", port)
	s.conn, _ = grpc.Dial(serverAddr, opts...)
	s.client = grpcserver.NewCalendarServiceClient(s.conn)

	wr := os.Stdout
	logger, err := logger.NewLogger(wr, "info")
	if err != nil {
		log.Fatal(errors.Wrapf(err, "can't init logger"))
	}
	s.repo, err = db.NewDBEventRepo("pgx", dsn, logger)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "can't init repo"))
	}
	s.testEvent, err = entities.NewEvent(title, dt, dur, text, userid, tnotify)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "can't init event"))
	}
	s.AfterTest("", "")
}

func (s *Suite) AfterTest(_, _ string) {
	//headers := validMetadata
	//ctx := metadata.NewOutgoingContext(context.Background(), headers)
	e, err := s.repo.GetForPeriod(context.Background(), time.Now().AddDate(-1, 0, 0), time.Now().AddDate(1, 0, 0))
	require.Nil(s.T(), err)
	for _, event := range e {
		s.repo.DeleteByID(context.Background(), event.ID)
	}
}

func (s *Suite) TestIntegration_AddEvent() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcserver.AddEventRequest
		postCondition condition
		err           bool
	}{
		{
			name:         "good",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.AddEventRequest{
				Title:      title,
				Datetime:   dt,
				Duration:   dur,
				Text:       text,
				Timenotify: tnotify,
			},
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:    "duplicate",
			headers: validMetadata,
			preCondition: func() {
				_, err := s.repo.Add(context.Background(), *s.testEvent)
				require.Nil(s.T(), err)
			},
			request: &grpcserver.AddEventRequest{
				Title:      title,
				Datetime:   date,
				Duration:   dur,
				Text:       text,
				Timenotify: tnotify,
			},
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: true,
		},
		{
			name:         "unauth",
			headers:      invalidMetadata,
			preCondition: nil,
			request: &grpcserver.AddEventRequest{
				Title:      title,
				Datetime:   date,
				Duration:   dur,
				Text:       text,
				Timenotify: tnotify,
			},
			postCondition: nil,
			err:           true,
		},
		{
			name:         "bad Duration value",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.AddEventRequest{
				Title:      title,
				Datetime:   date,
				Duration:   "bad",
				Text:       text,
				Timenotify: tnotify,
			},
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: true,
		},
		{
			name:         "too long title",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.AddEventRequest{
				Title:      "this title length is bigger than 100 letters.......this title length is bigger than 100 letters.......",
				Datetime:   date,
				Duration:   dur,
				Text:       text,
				Timenotify: tnotify,
			},
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: true,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.AddEvent(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
				require.Nil(s.T(), response)
			} else {
				require.Nil(s.T(), err)
				require.NotNil(s.T(), response)
				event, err := s.repo.GetByID(context.Background(), userid, response.GetId())
				require.Nil(s.T(), err)
				require.Equal(s.T(), s.testEvent.DateTime, event.DateTime)
				require.Equal(s.T(), s.testEvent.Title, event.Title)
				require.Equal(s.T(), s.testEvent.Text, event.Text)
				require.Equal(s.T(), s.testEvent.UserID, event.UserID)
				require.Equal(s.T(), s.testEvent.Duration, event.Duration)
			}
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}

func (s *Suite) TestIntegration_GetDate() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcserver.GetDateEventRequest
		responselen   int
		postCondition condition
		err           bool
	}{
		{
			name:    "good",
			headers: validMetadata,
			preCondition: func() {
				_, err := s.repo.Add(context.Background(), *s.testEvent)
				require.Nil(s.T(), err)
				// add 1 event after this day
				d, err := time.Parse(entities.LayoutISO, dt)
				require.Nil(s.T(), err)
				e, err := entities.NewEvent(title, d.AddDate(0, 0, 1).Format(entities.LayoutISO), dur, text, userid, tnotify)
				require.Nil(s.T(), err)
				_, err = s.repo.Add(context.Background(), *e)
				require.Nil(s.T(), err)

			},
			request: &grpcserver.GetDateEventRequest{Date: date,
			},
			responselen: 1,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:         "unauth",
			headers:      invalidMetadata,
			preCondition: nil,
			request: &grpcserver.GetDateEventRequest{Date: date,
			},
			responselen:   0,
			postCondition: nil,
			err:           true,
		},
		{
			name:         "not found",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.GetDateEventRequest{Date: date,
			},
			responselen: 0,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.GetDateEvent(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
				require.Nil(s.T(), response)
			} else {
				require.Nil(s.T(), err)
				require.NotNil(s.T(), response)
				events := response.GetEvents()
				require.Len(s.T(), events.Event, tcase.responselen)
			}
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}

func (s *Suite) TestIntegration_GetWeek() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcserver.GetWeekEventRequest
		responselen   int
		postCondition condition
		err           bool
	}{
		{
			name:    "good",
			headers: validMetadata,
			preCondition: func() {
				var e *entities.Event
				// add week events
				for i := 0; i < 7; i++ {
					d, err := time.Parse(entities.LayoutISO, dt)
					require.Nil(s.T(), err)
					e, err = entities.NewEvent(title, d.AddDate(0, 0, i).Format(entities.LayoutISO), dur, text, userid, tnotify)
					require.Nil(s.T(), err)
					_, err = s.repo.Add(context.Background(), *e)
					require.Nil(s.T(), err)
				}
				// add 3 events after this week
				for i := 7; i < 10; i++ {
					d, err := time.Parse(entities.LayoutISO, dt)
					require.Nil(s.T(), err)
					e, err = entities.NewEvent(title, d.AddDate(0, 0, i).Format(entities.LayoutISO), dur, text, userid, tnotify)
					require.Nil(s.T(), err)
					_, err = s.repo.Add(context.Background(), *e)
					require.Nil(s.T(), err)
				}
			},
			request: &grpcserver.GetWeekEventRequest{Date: date,
			},
			responselen: 7,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:         "unauth",
			headers:      invalidMetadata,
			preCondition: nil,
			request: &grpcserver.GetWeekEventRequest{Date: date,
			},
			responselen: 0,
			postCondition: nil,
			err:           true,
		},
		{
			name:         "not found",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.GetWeekEventRequest{Date: date,
			},
			responselen: 0,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.GetWeekEvent(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
				require.Nil(s.T(), response)
			} else {
				require.Nil(s.T(), err)
				require.NotNil(s.T(), response)
				events := response.GetEvents()
				require.Len(s.T(), events.Event, tcase.responselen)
			}
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}

func (s *Suite) TestIntegration_GetMonth() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcserver.GetMonthEventRequest
		responselen   int
		postCondition condition
		err           bool
	}{
		{
			name:    "good",
			headers: validMetadata,
			preCondition: func() {
				var e *entities.Event
				// add month events
				for i := 0; i < 31; i++ {
					d, err := time.Parse(entities.LayoutISO, dt)
					require.Nil(s.T(), err)
					e, err = entities.NewEvent(title, d.AddDate(0, 0, i).Format(entities.LayoutISO), dur, text, userid, tnotify)
					require.Nil(s.T(), err)
					_, err = s.repo.Add(context.Background(), *e)
					require.Nil(s.T(), err)
				}
				// add 5 events after this month
				for i := 31; i < 35; i++ {
					d, err := time.Parse(entities.LayoutISO, dt)
					require.Nil(s.T(), err)
					e, err = entities.NewEvent(title, d.AddDate(0, 0, i).Format(entities.LayoutISO), dur, text, userid, tnotify)
					require.Nil(s.T(), err)
					_, err = s.repo.Add(context.Background(), *e)
					require.Nil(s.T(), err)
				}
			},
			request: &grpcserver.GetMonthEventRequest{Date: date,
			},
			responselen: 31,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:         "unauth",
			headers:      invalidMetadata,
			preCondition: nil,
			request: &grpcserver.GetMonthEventRequest{Date: date,
			},
			responselen: 0,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err:           true,
		},
		{
			name:         "not found",
			headers:      validMetadata,
			preCondition: nil,
			request: &grpcserver.GetMonthEventRequest{Date: date,
			},
			responselen: 0,
			postCondition: func() {
				s.AfterTest("", "")
			},
			err: false,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.GetMonthEvent(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
				require.Nil(s.T(), response)
			} else {
				require.Nil(s.T(), err)
				require.NotNil(s.T(), response)
				events := response.GetEvents()
				require.Len(s.T(), events.Event, tcase.responselen)
			}
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}

/*
func (s *Suite) TestIntegration_ClickOnBanner() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcservice.ClickRequest
		postCondition condition
		err           bool
	}{
		{
			name:         "bad: banner not found",
			headers:      validMetadata,
			preCondition: func() { s.AfterTest("", "") },
			request: &grpcservice.ClickRequest{
				SlotId:   slotID,
				BannerId: 2,
				UserAge:  userAge,
				UserSex:  userSex,
			},
			postCondition: func() {
				group, err := s.repo.GetGroup(userAge, userSex)
				require.Nil(s.T(), err)
				// 10 tries - because queue makes delay
				clicks := uint(0)
				for i := 0; i < 10; i++ {
					actions, _ := s.repo.GetActions(pageURL, slotID, bannerID)
					if actions[*group].Clicks != 0 {
						clicks = actions[*group].Clicks
						break
					}
					time.Sleep(time.Second)
				}
				require.Equal(s.T(), uint(0), clicks)
				s.AfterTest("", "")
			},
			err: true,
		},
		{
			name:    "good",
			headers: validMetadata,
			preCondition: func() {
				err := s.repo.AddSlot(pageURL, slotID, slotDescription)
				require.Nil(s.T(), err)
				err = s.repo.AddBannerToSlot(pageURL, slotID, bannerID, bannerDescription)
				require.Nil(s.T(), err)
			},
			request: &grpcservice.ClickRequest{
				SlotId:   slotID,
				BannerId: bannerID,
				UserAge:  userAge,
				UserSex:  userSex,
			},
			postCondition: func() {
				group, err := s.repo.GetGroup(userAge, userSex)
				require.Nil(s.T(), err)
				// 10 tries - because queue makes delay
				clicks := uint(0)
				for i := 0; i < 10; i++ {
					actions, err := s.repo.GetActions(pageURL, slotID, bannerID)
					require.Nil(s.T(), err)
					if actions[*group].Clicks != 0 {
						clicks = actions[*group].Clicks
						break
					}
					time.Sleep(time.Second)
				}
				require.Equal(s.T(), uint(1), clicks)
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:    "unauth",
			headers: invalidMetadata,
			request: &grpcservice.ClickRequest{
				SlotId:   slotID,
				BannerId: bannerID,
				UserAge:  userAge,
				UserSex:  userSex,
			},
			err: true,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.ClickEvent(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
				require.Nil(s.T(), response)
			} else {
				require.Nil(s.T(), err)
				require.NotNil(s.T(), response)
			}
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}

func (s *Suite) TestIntegration_GetNextBanner() {
	tcases := []struct {
		name          string
		headers       metadata.MD
		preCondition  condition
		request       *grpcservice.GetNextBannerRequest
		response      *grpcservice.GetNextBannerResponse
		postCondition condition
		err           bool
	}{
		{
			name:         "bad: banner not found",
			headers:      validMetadata,
			preCondition: func() { s.AfterTest("", "") },
			request: &grpcservice.GetNextBannerRequest{
				SlotId:  2,
				UserAge: userAge,
				UserSex: userSex,
			},
			response: &grpcservice.GetNextBannerResponse{
				BannerId: 0,
			},
			postCondition: func() {
				group, err := s.repo.GetGroup(userAge, userSex)
				require.Nil(s.T(), err)
				// 10 tries - because queue makes delay
				shows := uint(0)
				for i := 0; i < 10; i++ {
					actions, _ := s.repo.GetActions(pageURL, slotID, bannerID)
					if actions[*group].Shows != 0 {
						shows = actions[*group].Shows
						break
					}
					time.Sleep(time.Second)
				}
				require.Equal(s.T(), uint(0), shows)
				s.AfterTest("", "")
			},
			err: true,
		},
		{
			name:    "good",
			headers: validMetadata,
			preCondition: func() {
				err := s.repo.AddSlot(pageURL, slotID, slotDescription)
				require.Nil(s.T(), err)
				err = s.repo.AddBannerToSlot(pageURL, slotID, bannerID, bannerDescription)
				require.Nil(s.T(), err)
			},
			request: &grpcservice.GetNextBannerRequest{
				SlotId:  slotID,
				UserAge: userAge,
				UserSex: userSex,
			},
			response: &grpcservice.GetNextBannerResponse{
				BannerId: bannerID,
			},
			postCondition: func() {
				group, err := s.repo.GetGroup(userAge, userSex)
				require.Nil(s.T(), err)
				// 10 tries - because queue makes delay
				shows := uint(0)
				for i := 0; i < 10; i++ {
					actions, err := s.repo.GetActions(pageURL, slotID, bannerID)
					require.Nil(s.T(), err)
					if actions[*group].Shows != 0 {
						shows = actions[*group].Shows
						break
					}
					time.Sleep(time.Second)
				}
				require.Equal(s.T(), uint(1), shows)
				s.AfterTest("", "")
			},
			err: false,
		},
		{
			name:    "unauth",
			headers: invalidMetadata,
			request: &grpcservice.GetNextBannerRequest{
				SlotId:  slotID,
				UserAge: userAge,
				UserSex: userSex,
			},
			err: true,
		},
	}

	for id, tcase := range tcases {
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			if tcase.preCondition != nil {
				tcase.preCondition()
			}
			ctx := metadata.NewOutgoingContext(context.Background(), tcase.headers)
			response, err := s.client.GetNextBanner(ctx, tcase.request)
			if tcase.err {
				require.NotNil(s.T(), err)
			} else {
				require.Nil(s.T(), err)
			}
			require.Equal(s.T(), tcase.response.GetBannerId(), response.GetBannerId())
			if tcase.postCondition != nil {
				tcase.postCondition()
			}
		})
	}
}*/
