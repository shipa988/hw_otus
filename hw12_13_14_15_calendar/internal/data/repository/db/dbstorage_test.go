package db

import (
	"context"
	"database/sql"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

var (
	testid     = "00112233-4455-6677-8899-aabbccddeeff"
	testid2    = "11112233-4455-6677-8899-aabbccddeeff"
	testid3    = "22112233-4455-6677-8899-aabbccddeeff"
	testdt_str = "2020-06-25 14:14:14"
	testdt, _  = time.Parse(entities.LayoutISO, testdt_str)
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
			AddRow(testid, "title", testdt, 360, "text", testid, 360)

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify 
												from`).
			WithArgs(testid, testid).
			WillReturnRows(rows)

		expectedEvent := addID(newFakeEvent(), testid)

		event, err := dbe.GetByID(ctx, testid, testid)

		if mockerr := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Nil(t, err)
		require.Equal(t, &expectedEvent, event)
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
			WithArgs(testid, testid).
			WillReturnRows(rows)

		dbe := EventRepo{db: db, logger: nil}
		event, err := dbe.GetByID(ctx, testid, testid)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Nil(t, event, "no events should be returned")
		require.Equal(t, entities.ErrEventNotFound, err, "should return event not found")
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
		event, err := dbe.GetByID(ctx, testid, testid)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Nil(t, event, "no events should be returned")
		require.Truef(t, errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")
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
			AddRow(testid, "title", testdt, 360, "text", testid, 360)

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(testdt, testdt, testid).
			WillReturnRows(rows)

		var expectedEvents []*entities.Event
		expectedEvent := addID(newFakeEvent(), testid)
		expectedEvents = append(expectedEvents, &expectedEvent)

		events, err := dbe.GetForPeriod(ctx, testid, testdt, testdt)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equal(t, expectedEvents, events)
		require.Nil(t, err)
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
			WithArgs(testdt, testdt, testid).
			WillReturnRows(rows)

		dbe := EventRepo{db: db, logger: nil}
		var expectedEvents []*entities.Event //nil

		events, err := dbe.GetForPeriod(ctx, testid, testdt, testdt)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equalf(t, expectedEvents, events, "")
		require.Equal(t, entities.ErrEventNotFound, err, "should return event not found")

	})
	t.Run("return error: get event for period", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.TODO()

		mock.ExpectQuery(`select id, title, datetime, duration, text, userid, timenotify from`).
			WithArgs(testdt, testdt, testid).
			WillReturnError(sql.ErrConnDone)

		dbe := EventRepo{db: db, logger: nil}
		var expectedEvents []*entities.Event //nil

		events, err := dbe.GetForPeriod(ctx, testid, testdt, testdt)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equalf(t, expectedEvents, events, "")
		require.Truef(t, errors.Is(err, sql.ErrConnDone), "GetForPeriod not return cause error")

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
			WithArgs("title", testdt, 360, "text", testid, 360).
			WillReturnRows(rows)

		event := newFakeEvent()
		id, err := dbe.Add(ctx, event)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equal(t, testid, id, "generated id must be non empty")
		require.Nil(t, err)
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
			WithArgs("title", testdt, 360, "text", testid, 360).
			WillReturnError(sql.ErrConnDone)

		event := newFakeEvent()
		id, err := dbe.Add(ctx, event)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equal(t, "", id, "generated id must be empty")
		require.True(t, errors.Is(err, sql.ErrConnDone))
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
			WithArgs("title", testdt, 360, "text", testid, 360).
			WillReturnError(errors.New(uniqueViolation))

		event := newFakeEvent()
		id, err := dbe.Add(ctx, event)

		if mockerr := mock.ExpectationsWereMet(); mockerr != nil {
			t.Errorf("there were unfulfilled expectations: %s", mockerr)
		}

		require.Equal(t, "", id, "generated id must be empty")
		require.True(t, errors.Is(err, entities.ErrDateBusy))
		require.Containsf(t, err.Error(), entities.ErrDateBusy.Error(), "")
	})
}
func TestEventConversion(t *testing.T) {
	t.Run("DB Event to Domain Event conversion", func(t *testing.T) {
		dbE := Event{
			ID:         uuid.FromStringOrNil(testid),
			Title:      "title",
			DateTime:   testdt,
			Duration:   360,
			Text:       "text",
			UserID:     uuid.FromStringOrNil(testid),
			TimeNotify: 360,
		}
		expectedE := addID(newFakeEvent(), testid)
		e := toDomainEvent(dbE)
		require.Equal(t, &expectedE, e, "domain event not equal db event after conversion")
	})
	t.Run("event to DB Event conversion", func(t *testing.T) {
		event := addID(newFakeEvent(), testid)
		expectedDBE := &Event{
			ID:         uuid.FromStringOrNil(testid),
			Title:      "title",
			DateTime:   testdt,
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

func newFakeEvent() entities.Event {
	e, _ := entities.NewEvent("title", testdt_str, "6m", "text", testid, "6m")
	return *e
}
func addID(e entities.Event, id string) entities.Event {
	e.ID = id
	return e
}
