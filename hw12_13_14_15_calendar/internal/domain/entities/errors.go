package entities

import (
	"github.com/pkg/errors"
)

// ErrorEntity is business error.
var (
	ErrDateBusy       = errors.New("date for this event is busy")
	ErrUnknownUser    = errors.New("user not found")
	ErrEventNotFound  = errors.New("event not found")
	ErrTitleLen       = errors.New("event title field is to long. It should be no more than 100 letters")
	ErrDateFormat     = errors.New("date of event is in the incorrect format. format is " + LayoutDateISO)
	ErrDatetimeFormat = errors.New("date of event is in the incorrect format. format is " + LayoutISO)
	ErrDurationFormat = errors.New("duration of the event is in the incorrect format. format is: XXs,XXm,XXh,XXd, where XX-decimal")
	ErrNotifyFormat   = errors.New("time to notify about event is in the incorrect format. format is: XXs,XXm,XXh,XXd, where XX-decimal")
	ErrEventFormat    = errors.New("not correct event")
)

var ErrNoField = "field %v is necessary"
