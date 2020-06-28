package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/usecases"
	"github.com/stretchr/testify/require"
)

var (
	testid     = "00112233-4455-6677-8899-aabbccddeeff"
	testid2    = "11112233-4455-6677-8899-aabbccddeeff"
	testid3    = "22112233-4455-6677-8899-aabbccddeeff"
	testdt_str = "2020-06-25 14:14:14"
	testdt, _  = time.Parse(usecases.LayoutISO, testdt_str)
)

func TestDBEventRepo_GetById(t *testing.T) {
	t.Run("good test: get event by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbe := EventRepo{db: db, logger: nil}

		ctx := context.TODO()

		rows := sqlmock.NewRows([]string{"id", "title", "dateTime", "duration", "text", "userId", "timeNotify"}).
			AddRow(testid, "title", testdt_str, 360, "text", testid, 360)

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify 
												from`).
			WithArgs(testid).
			WillReturnRows(rows)

		expectedEvent := newFakeEvent(testid)

		event, err := dbe.GetByID(ctx, testid)
		require.Nil(t, err)
		require.Equal(t, &expectedEvent, event)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("no rows: get event by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.TODO()

		rows := sqlmock.NewRows([]string{"id", "title", "dateTime", "duration", "text", "userId", "timeNotify"})
		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(testid).
			WillReturnRows(rows)

		dbe := EventRepo{db: db, logger: nil}
		event, err := dbe.GetByID(ctx, testid)

		require.Nil(t, event, "no events should be returned")
		require.Equal(t, domain.ErrEventNotFound, err, "should return event not found")

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error: get event by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.TODO()

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WillReturnError(sql.ErrConnDone)

		dbe := EventRepo{db: db, logger: nil}
		event, err := dbe.GetByID(ctx, testid)

		require.Nil(t, event, "no events should be returned")
		require.Truef(t, errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestDBEventRepo_GetForPeriod(t *testing.T) {
	t.Run("good test: get event for period", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbe := EventRepo{db: db, logger: nil}

		ctx := context.TODO()

		rows := sqlmock.NewRows([]string{"id", "title", "dateTime", "duration", "text", "userId", "timeNotify"}).
			AddRow(testid, "title", testdt_str, 360, "text", testid, 360)

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(time.Now(), time.Now()).
			WillReturnRows(rows)

		var expectedEvents []*domain.Event
		expectedEvent := newFakeEvent(testid)
		expectedEvents = append(expectedEvents, &expectedEvent)

		events, err := dbe.GetForPeriod(ctx, time.Now(), time.Now())
		require.Equal(t, expectedEvents, events)
		require.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("no rows: get event for period", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.TODO()

		rows := sqlmock.NewRows([]string{"id", "title", "dateTime", "duration", "text", "userId", "timeNotify"})
		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(time.Now(), time.Now()).
			WillReturnRows(rows)

		dbe := EventRepo{db: db, logger: nil}
		var expectedEvents []*domain.Event //nil

		events, err := dbe.GetForPeriod(ctx, time.Now(), time.Now())

		require.Equalf(t, expectedEvents, events, "")
		require.Equal(t, domain.ErrEventNotFound, err, "should return event not found")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error: get event for period", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.TODO()

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(time.Now(), time.Now()).
			WillReturnError(sql.ErrConnDone)

		dbe := EventRepo{db: db, logger: nil}
		var expectedEvents []*domain.Event //nil

		events, err := dbe.GetForPeriod(ctx, time.Now(), time.Now())
		require.Equalf(t, expectedEvents, events, "")
		require.Truef(t, errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestDBEventRepo_Add(t *testing.T) {
	t.Run("good test: add event", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbe := EventRepo{db: db, logger: nil}

		ctx := context.TODO()
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(testid)

		mock.ExpectQuery(`INSERT INTO public.events`).
			WithArgs("title", testdt.String(), float64(360), "text", testid, float64(360)).
			WillReturnRows(rows)

		event := domain.NewEvent("title", testdt, 360000000000, "text", testid, 360000000000)

		id, err := dbe.Add(ctx, *event)
		require.Equal(t, testid, id, "generated id must be non empty")
		require.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error: add event", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbe := EventRepo{db: db, logger: nil}

		ctx := context.TODO()

		mock.ExpectQuery(`INSERT INTO public.events`).
			WithArgs("title", testdt.String(), float64(360), "text", testid, float64(360)).
			WillReturnError(sql.ErrConnDone)

		event := domain.NewEvent("title", testdt, 360000000000, "text", testid, 360000000000)

		id, err := dbe.Add(ctx, *event)
		require.Equal(t, "", id, "generated id must be empty")
		require.True(t, errors.Is(err, sql.ErrConnDone))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
	t.Run("return date is busy error: add event", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbe := EventRepo{db: db, logger: nil}

		ctx := context.TODO()

		mock.ExpectQuery(`INSERT INTO public.events`).
			WithArgs("title", testdt.String(), float64(360), "text", testid, float64(360)).
			WillReturnError(errors.New(UniqueViolation))

		event := domain.NewEvent("title", testdt, 360000000000, "text", testid, 360000000000)

		id, err := dbe.Add(ctx, *event)
		require.Equal(t, "", id, "generated id must be empty")
		require.True(t, errors.Is(err, domain.ErrDateBusy))
		require.Containsf(t, err.Error(), domain.ErrDateBusy.Error(), "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestEventConversion(t *testing.T) {
	t.Run("DB Event to Domain Event conversion", func(t *testing.T) {
		dbE := Event{
			ID:         uuid.FromStringOrNil(testid),
			Title:      "title",
			DateTime:   testdt_str,
			Duration:   360,
			Text:       "text",
			UserID:     uuid.FromStringOrNil(testid),
			TimeNotify: 360,
		}
		expectedE := newFakeEvent(testid)
		e, err := toDomainEvent(dbE)
		require.Nil(t, err, "error must be nil")
		require.Equal(t, &expectedE, e, "domain event not equal db event after conversion")
	})
	t.Run("event to DB Event conversion", func(t *testing.T) {
		event := newFakeEvent(testid)
		expectedDBE := &Event{
			ID:         uuid.FromStringOrNil(testid),
			Title:      "title",
			DateTime:   testdt.String(),
			Duration:   360,
			Text:       "text",
			UserID:     uuid.FromStringOrNil(testid),
			TimeNotify: 360,
		}
		dbe, err := fromDomainEvent(event)
		require.Nil(t, err, "error must be nil")
		require.Equal(t, expectedDBE, dbe, "domain event not equal db event after conversion")
	})
}

func newFakeEvent(id string) domain.Event {
	e := domain.NewEvent("title", testdt, 360000000000, "text", id, 360000000000)
	e.ID = id
	return *e
}
