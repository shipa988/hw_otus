package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	LayoutISO     = "2006-01-02 15:04:05"
	LayoutDateISO = "2006-01-02"
)

type EventRepo interface {
	Add(ctx context.Context, event Event) (ID string, err error)
	GetByID(ctx context.Context, userID, eventID string) (*Event, error)
	GetByDate(ctx context.Context, userID string, date time.Time) ([]*Event, error)
	GetForPeriod(ctx context.Context, userID string, dateStart time.Time, dateEnd time.Time) ([]*Event, error)
	UpdateByID(ctx context.Context, userID, eventID string, event Event) error
	DeleteByID(ctx context.Context, userID, eventID string) error
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

func (e Event) Update(title, dateTime, duration, text, timeNotify string) (Event, error) {
	//update duration if not empty
	if duration != "" {
		dr, err := time.ParseDuration(duration)
		if err != nil {
			return e, errors.Wrapf(ErrDurationFormat, "error update duration: %v", duration)
		}
		e.Duration = dr
	}
	//update timeNotify if not empty
	if timeNotify != "" {
		tn, err := time.ParseDuration(timeNotify)
		if err != nil {
			return e, errors.Wrapf(ErrNotifyFormat, "error update timeNotify: %v", timeNotify)
		}
		e.TimeNotify = tn
	}
	//update dateTime if not empty
	if dateTime != "" {
		dt, err := time.Parse(LayoutISO, dateTime)
		if err != nil {
			return e, errors.Wrapf(ErrDatetimeFormat, "error update dateTime: %v", dateTime)
		}
		e.DateTime = dt
	}
	//update title
	if title != "" && title != e.Title {
		e.Title = title
	}
	//update text
	if text != "" && text != e.Text {
		e.Text = text
	}
	return e, nil
}

// NewEvent create Event entity and validating business logic of event.
func NewEvent(title, dateTime, duration, text, userID, timeNotify string) (*Event, error) {
	if title == "" {
		return nil, fmt.Errorf(ErrNoField, "title")
	}
	if dateTime == "" {
		return nil, fmt.Errorf(ErrNoField, "dateTime")
	}
	if text == "" {
		return nil, fmt.Errorf(ErrNoField, "text")
	}
	if userID == "" {
		return nil, fmt.Errorf(ErrNoField, "userID")
	}
	if duration == "" {
		return nil, fmt.Errorf(ErrNoField, "duration")
	}

	if len(title) > 100 {
		return nil, ErrTitleLen
	}
	tl := title

	dt, err := time.Parse(LayoutISO, dateTime)
	if err != nil {
		return nil, errors.Wrapf(ErrDatetimeFormat, "error set dateTime: %v", dateTime)
	}

	dr, err := time.ParseDuration(duration)
	if err != nil {
		return nil, errors.Wrapf(ErrDurationFormat, "error set duration: %v", duration)
	}

	tn := time.Duration(0)
	if timeNotify != "" {
		tn, err = time.ParseDuration(timeNotify)
		if err != nil {
			return nil, errors.Wrapf(ErrNotifyFormat, "error set timeNotify:  %v", timeNotify)
		}
	}

	return &Event{
		Title:      tl,
		DateTime:   dt,
		Duration:   dr,
		Text:       text,
		UserID:     userID,
		TimeNotify: tn,
	}, nil
}
