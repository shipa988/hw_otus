package domain

import "github.com/pkg/errors"

// ErrorEntity is business error.
var (
	ErrDateBusy      = errors.New("date for this event is busy")
	ErrEventNotFound = errors.New("event not found")
	ErrTitleLen      = errors.New("event title field is to long. It should be no more than 100 letters")
	ErrEventFormat   = errors.New("not correct event")
)
