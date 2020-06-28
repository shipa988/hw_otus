package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/usecases"
)

//infrastructure layer error.
const (
	UniqueViolation = "23505" //UniqueViolation              = "23505"
)

//interface layer error.
const (
	ErrAdd          = "can't add to database event: %v"
	ErrGetbyID      = "can't get event from database by id: %v"
	ErrGetbyDate    = "can't get event from database by date: %v"
	ErrGetForPeriod = "can't get event from database for period: %v-%v"
	ErrUpdatebyID   = "can't update event in database by id: %v"
	ErrDeletebyID   = "can't delete event in database by id: %v"
	ErrConvert      = "can't convert between business event and database event entities"
)

var _ domain.EventRepo = (*EventRepo)(nil)

type EventRepo struct {
	db     *sql.DB
	logger usecases.ILogger
}

type Event struct {
	ID         uuid.UUID
	Title      string
	DateTime   string
	Duration   float64
	Text       string
	UserID     uuid.UUID
	TimeNotify float64
}

func NewDBEventRepo(driver, dsn string, logger usecases.ILogger) (*EventRepo, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create connect to db with dsn %v by driver %v", dsn, driver)
	}
	return &EventRepo{
		db:     db,
		logger: logger,
	}, nil
}

func (repo *EventRepo) Add(ctx context.Context, event domain.Event) (string, error) {
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
		if strings.Contains(err.Error(), UniqueViolation) {
			//yes, it is bad practice but we are inside the
			//"interface layer" and we could't use driver errors such as https://github.com/jackc/pgx/wiki/Error-Handling,
			//because driver - "infrastructure layer" (dependency direction1)
			return "", domain.ErrDateBusy
		}
		return "", errors.Wrapf(err, ErrAdd, event)
	}
	return id, nil
}

func (repo *EventRepo) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	row := repo.db.QueryRowContext(ctx, `select id, title, datetime, duration, text, userid, timenotify 
												from public.events where id = $1`, id)
	if row == nil {
		return nil, domain.ErrEventNotFound
	}
	dbevent := Event{}
	err := row.Scan(&dbevent.ID, &dbevent.Title, &dbevent.DateTime, &dbevent.Duration, &dbevent.Text, &dbevent.UserID, &dbevent.TimeNotify)
	if err != nil {
		return nil, SQLError(err, fmt.Sprintf(ErrGetbyID, id))
	}
	event, err := toDomainEvent(dbevent)
	if err != nil {
		return nil, errors.Wrap(err, ErrConvert)
	}
	return event, nil
}

func (repo *EventRepo) GetByDate(ctx context.Context, date time.Time) ([]*domain.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify 
												 	from public.events where cast (datetime as date)=cast ($1 as date) =$1`, date)
	if err != nil {
		return nil, SQLError(err, fmt.Sprintf(ErrGetbyDate, date))
	}
	defer rows.Close()

	events := []*domain.Event{}
	for rows.Next() {
		dbevent := Event{}
		err := rows.Scan(&dbevent.ID, &dbevent.Title, &dbevent.DateTime, &dbevent.Duration, &dbevent.Text, &dbevent.UserID, &dbevent.TimeNotify)
		if err != nil {
			return nil, SQLError(err, fmt.Sprintf(ErrGetbyDate, date))
		}
		event, err := toDomainEvent(dbevent)
		if err != nil {
			return nil, errors.Wrap(err, ErrConvert)
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, domain.ErrEventNotFound
	}

	if rows.Err() != nil && err != sql.ErrNoRows {
		return nil, errors.Wrapf(rows.Err(), ErrGetbyDate, date)
	}
	return events, nil
}

func (repo *EventRepo) GetForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]*domain.Event, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, title, datetime, duration, text, userid, timenotify from public.events where datetime between  $1 and $2`, dateStart, dateEnd)
	if err != nil {
		return nil, SQLError(err, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
	}
	defer rows.Close()

	events := []*domain.Event{}
	for rows.Next() {
		dbevent := Event{}
		err := rows.Scan(&dbevent.ID, &dbevent.Title, &dbevent.DateTime, &dbevent.Duration, &dbevent.Text, &dbevent.UserID, &dbevent.TimeNotify)
		if err != nil {
			return nil, SQLError(err, fmt.Sprintf(ErrGetForPeriod, dateStart, dateEnd))
		}
		event, err := toDomainEvent(dbevent)
		if err != nil {
			return nil, errors.Wrap(err, ErrConvert)
		}
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, domain.ErrEventNotFound
	}

	if rows.Err() != nil && err != sql.ErrNoRows {
		return nil, errors.Wrapf(rows.Err(), ErrGetForPeriod, dateStart, dateEnd)
	}
	return events, nil
}

func (repo *EventRepo) UpdateByID(ctx context.Context, id string, event domain.Event) error {
	dbEvent, err := fromDomainEvent(event)
	if err != nil {
		return errors.Wrap(err, ErrConvert)
	}

	result, err := repo.db.ExecContext(ctx, `UPDATE public.events
	SET title=$2, datetime=$3, duration=$4, text=$5, userid=$6, timenotify=$7
	WHERE id=$1;`, id, dbEvent.Title, dbEvent.DateTime, dbEvent.Duration, dbEvent.Text, dbEvent.UserID, dbEvent.TimeNotify)

	if err != nil {
		return errors.Wrapf(err, ErrUpdatebyID, id)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrUpdatebyID, id)
	}
	if rows != 1 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (repo *EventRepo) DeleteByID(ctx context.Context, id string) error {
	result, err := repo.db.ExecContext(ctx, `DELETE FROM public.events where id = $1`, id)
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, id)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrDeletebyID, id)
	}
	if rows != 1 {
		return domain.ErrEventNotFound
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

func toDomainEvent(dbe Event) (*domain.Event, error) {
	dt, err := time.Parse(usecases.LayoutISO, dbe.DateTime)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert entities DB Event->Event ")
	}
	event := &domain.Event{
		ID:         dbe.ID.String(),
		Title:      dbe.Title,
		DateTime:   dt,
		Duration:   time.Second * time.Duration(dbe.Duration),
		Text:       dbe.Text,
		UserID:     dbe.UserID.String(),
		TimeNotify: time.Second * time.Duration(dbe.TimeNotify),
	}
	return event, nil
}
func fromDomainEvent(event domain.Event) (*Event, error) {
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
	dbe.Duration = event.Duration.Seconds()
	dbe.TimeNotify = event.TimeNotify.Seconds()
	dbe.UserID = uid
	dbe.DateTime = event.DateTime.String()
	return dbe, nil
}

func SQLError(err error, message string) error {
	switch err {
	case sql.ErrNoRows:
		return domain.ErrEventNotFound
	default:
		return errors.Wrap(err, message)
	}
}
