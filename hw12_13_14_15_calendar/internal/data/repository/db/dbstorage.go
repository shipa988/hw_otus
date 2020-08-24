package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

//infrastructure layer error.
const (
	uniqueViolation = "23505" //uniqueViolation              = "23505"
)

//interface layer error.
const (
	ErrAdd             = "can't add to database event: %v"
	ErrGetbyID         = "can't get event from database by id: %v"
	ErrGetbyDate       = "can't get event from database by date: %v"
	ErrGetbyNotifyDate = "can't get event from database by notify date: %v"
	ErrGetForPeriod    = "can't get event from database for period: %v-%v"
	ErrUpdatebyID      = "can't update event in database by id: %v"
	ErrDeletebyID      = "can't delete event in database by id: %v"
	ErrConvert         = "can't convert between business event and database event entities"
)

var _ entities.EventRepo = (*EventRepo)(nil)

type EventRepo struct {
	db     *sql.DB
	logger usecases.Logger
}

type Event struct {
	ID         uuid.UUID
	Title      string
	DateTime   time.Time
	Duration   int
	Text       string
	UserID     uuid.UUID
	TimeNotify int
}

func NewDBEventRepo(driver, dsn string, logger usecases.Logger) (*EventRepo, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create connect to db with dsn %v by driver %v", dsn, driver)
	}
	return &EventRepo{
		db:     db,
		logger: logger,
	}, nil
}

func (repo *EventRepo) Add(ctx context.Context, event entities.Event) (string, error) {
	dbEvent, err := fromDomainEvent(event)
	if err != nil {
		return "", errors.Wrap(err, ErrConvert)
	}
	var id string

	row := repo.db.QueryRowContext(ctx, `INSERT INTO public.events(
	id, title, datetime, duration, text, userid, timenotify)
	values (uuid_generate_v4(),$1, $2, $3, $4, $5, $6) returning id`, dbEvent.Title, dbEvent.DateTime, dbEvent.Duration, dbEvent.Text, dbEvent.UserID, dbEvent.TimeNotify)

	if row == nil {
		return "", errors.Wrapf(errors.New("insert query return nil row"), ErrAdd, event)
	}

	err = row.Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), uniqueViolation) {
			//yes, it is bad practice but we are inside the
			//"interface layer" and we could't use driver errors such as https://github.com/jackc/pgx/wiki/Error-Handling,
			//because driver - "infrastructure layer" (dependency direction1)
			return "", entities.ErrDateBusy
		}
		return "", errors.Wrapf(err, ErrAdd, event)
	}
	return id, nil
}

func (repo *EventRepo) GetByID(ctx context.Context, userID, eventID string) (*entities.Event, error) {
	row := repo.db.QueryRowContext(ctx, `select id, title, datetime, duration, text, userid, timenotify 
												from public.events where userid=$1 and id = $2`, userID, eventID)
	if row == nil {
		return nil, entities.ErrEventNotFound
	}
	dbevent := Event{}
	err := row.Scan(&dbevent.ID, &dbevent.Title, &dbevent.DateTime, &dbevent.Duration, &dbevent.Text, &dbevent.UserID, &dbevent.TimeNotify)
	if err != nil {
		return nil, SQLError(err, fmt.Sprintf(ErrGetbyID, eventID))
	}
	event := toDomainEvent(dbevent)
	return event, nil
}

func (repo *EventRepo) GetByDate(ctx context.Context, userID string, date time.Time) ([]*entities.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify 
												 	from public.events where cast (datetime as date)=cast ($1 as date) and userid=$2`, date, userID)

	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, fmt.Sprintf(ErrGetbyDate, date))
	}
	defer rows.Close()

	return repo.rowsToEvents(rows, fmt.Sprintf(ErrGetbyDate, date))
}

func (repo *EventRepo) GetByNotifyDate(ctx context.Context, date time.Time) ([]*entities.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify from public.events where timenotify is not null and (cast (datetime as date) - make_interval(secs => timenotify))=cast ($1 as date)`, date)
	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, fmt.Sprintf(ErrGetbyNotifyDate, date))
	}
	defer rows.Close()

	return repo.rowsToEvents(rows, fmt.Sprintf(ErrGetbyNotifyDate, date))
}

func (repo *EventRepo) GetForPeriodByUserID(ctx context.Context, userID string, dateStart time.Time, dateEnd time.Time) ([]*entities.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify from public.events where (datetime between  $1 and $2)  and userid=$3`, dateStart, dateEnd, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
	}
	defer rows.Close()

	return repo.rowsToEvents(rows, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
}

func (repo *EventRepo) GetForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]*entities.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify from public.events where (datetime between  $1 and $2)`, dateStart, dateEnd)
	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
	}
	defer rows.Close()

	return repo.rowsToEvents(rows, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
}
func (repo *EventRepo) UpdateByID(ctx context.Context, userID, eventID string, event entities.Event) error {
	dbEvent, err := fromDomainEvent(event)
	if err != nil {
		return errors.Wrap(err, ErrConvert)
	}

	result, err := repo.db.ExecContext(ctx, `UPDATE public.events
	SET title=$3, datetime=$4, duration=$5, text=$6, userid=$7, timenotify=$8
	WHERE  userid=$1 and id=$2;`, userID, eventID, dbEvent.Title, dbEvent.DateTime, dbEvent.Duration, dbEvent.Text, dbEvent.UserID, dbEvent.TimeNotify)

	if err != nil {
		return errors.Wrapf(err, ErrUpdatebyID, eventID)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrUpdatebyID, eventID)
	}
	if rows != 1 {
		return entities.ErrEventNotFound
	}

	return nil
}

func (repo *EventRepo) DeleteByUserID(ctx context.Context, userID, eventID string) error {
	result, err := repo.db.ExecContext(ctx, `DELETE FROM public.events where  userid=$1 and id=$2;`, userID, eventID)
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, eventID)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, eventID)
	}
	if rows != 1 {
		return entities.ErrEventNotFound
	}

	return nil
}

func (repo *EventRepo) DeleteByID(ctx context.Context, eventID string) error {
	result, err := repo.db.ExecContext(ctx, `DELETE FROM public.events where id=$1;`, eventID)
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, eventID)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, eventID)
	}
	if rows != 1 {
		return entities.ErrEventNotFound
	}

	return nil
}

func (repo *EventRepo) Connect(ctx context.Context, dsn string) (err error) {
	err = repo.db.PingContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to db: %v", dsn)
	}
	return nil
}

func (repo *EventRepo) Close() error {
	return repo.db.Close()
}

func (repo *EventRepo) rowsToEvents(rows *sql.Rows, errorString string) ([]*entities.Event, error) {
	events := []*entities.Event{}
	for rows.Next() {
		dbevent := Event{}
		err := rows.Scan(&dbevent.ID, &dbevent.Title, &dbevent.DateTime, &dbevent.Duration, &dbevent.Text, &dbevent.UserID, &dbevent.TimeNotify)
		if err != nil {
			return nil, SQLError(err, errorString)
		}
		event := toDomainEvent(dbevent)
		events = append(events, event)
	}

	/*if len(events) == 0 {
		return nil, entities.ErrEventNotFound
	}*/

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), errorString)
	}
	return events, nil
}

func toDomainEvent(dbe Event) *entities.Event {
	event := &entities.Event{
		ID:         dbe.ID.String(),
		Title:      dbe.Title,
		DateTime:   dbe.DateTime,
		Duration:   time.Second * time.Duration(dbe.Duration),
		Text:       dbe.Text,
		UserID:     dbe.UserID.String(),
		TimeNotify: time.Second * time.Duration(dbe.TimeNotify),
	}
	return event
}
func fromDomainEvent(event entities.Event) (*Event, error) {
	id, err := uuid.FromString(event.ID)
	if err != nil && event.ID != "" {
		return nil, errors.Wrap(err, "can't convert domain event to db event")
	}
	uid, err := uuid.FromString(event.UserID)
	if err != nil && event.UserID != "" {
		return nil, errors.Wrap(err, "can't convert domain event to db event")
	}
	dbe := &Event{}
	dbe.ID = id
	dbe.Title = event.Title
	dbe.Text = event.Text
	dbe.Duration = int(event.Duration.Seconds())
	dbe.TimeNotify = int(event.TimeNotify.Seconds())
	dbe.UserID = uid
	dbe.DateTime = event.DateTime
	return dbe, nil
}

func SQLError(err error, message string) error {
	switch err {
	case sql.ErrNoRows:
		return entities.ErrEventNotFound
	default:
		return errors.Wrap(err, message)
	}
}
