package domain

import (
	"context"
	"time"
)

type EventRepo interface {
	Add(ctx context.Context, event Event) (ID string, err error)
	GetByID(ctx context.Context, id string) (*Event, error)
	GetByDate(ctx context.Context, date time.Time) ([]*Event, error)
	GetForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]*Event, error)
	UpdateByID(ctx context.Context, id string, event Event) error
	DeleteByID(ctx context.Context, id string) error
}

type Event struct {
	ID         string
	Title      string
	DateTime   time.Time
	Duration   time.Duration
	Text       string
	UserID     string
	TimeNotify time.Duration
}

func NewEvent(title string, dateTime time.Time, duration time.Duration, text string, userID string, timeNotify time.Duration) *Event {
	return &Event{
		Title:      title,
		DateTime:   dateTime,
		Duration:   duration,
		Text:       text,
		UserID:     userID,
		TimeNotify: timeNotify,
	}
}

// Validate validating business logic of event.
func (e Event) Validate() error {
	if len(e.Title) > 100 {
		e.Title = e.Title[:100]
		return ErrTitleLen
	}
	return nil
}

type NotifyRepo interface {
	Add(ctx context.Context, alert Notify) error
}

type Notify struct {
	ID       string
	Title    string
	UserID   string
	DateTime time.Time
}

type User struct {
	ID string
}

/*
type UserRepo interface {
	Create(User)
	Delete(uuid.UUID)
}

type Scheduler struct {
	EventID uuid.UUID
	UserID uuid.UUID
}

type Notifier struct {
	EventID uuid.UUID
	UserID uuid.UUID
}
*/
