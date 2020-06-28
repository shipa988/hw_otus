package infrastructure

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shipa988/otus/hw12_13_14_15_calendar/internal/usecases"
)

var _ usecases.ILogger = (*Logger)(nil)

type Logger struct {
	logger *zerolog.Logger
}

func NewLogger(logWriter io.Writer, logLevel string) (*Logger, error) {
	ll := zerolog.InfoLevel

	switch strings.ToLower(logLevel) {
	case "error":
		ll = zerolog.ErrorLevel
	case "warn":
		ll = zerolog.WarnLevel
	case "info":
		ll = zerolog.InfoLevel
	case "debug":
		ll = zerolog.DebugLevel
	default:
		log.Println("unknown log level, I know levels:error,warn,info,debug. Setting default level:info")
	}
	zerolog.SetGlobalLevel(ll)
	logger := zerolog.New(logWriter).With().Timestamp().Logger()
	return &Logger{logger: &logger}, nil
}

func (l *Logger) Debug(ctx context.Context, message interface{}, args ...interface{}) {
	switch mess := message.(type) {
	case error:
		err, ok := errors.Cause(mess).(stackTracer)
		if !ok {
			l.Debug(ctx, mess.Error(), args)
			return
		}
		st := err.StackTrace()
		l.debug(ctx, fmt.Sprintf("%+v", st), args)
	case string:
		l.debug(ctx, mess, args)
	default:
		l.debug(ctx, fmt.Sprintf("debug message %v has unknown type %v", message, mess), args)
	}
}
func (l *Logger) Error(ctx context.Context, message interface{}, args ...interface{}) {
	switch mess := message.(type) {
	case error:
		err, ok := errors.Cause(mess).(stackTracer)
		if !ok {
			l.Debug(ctx, mess.Error(), args)
			return
		}
		st := err.StackTrace()
		l.error(ctx, fmt.Sprintf("%v", st), args)
	case string:
		l.error(ctx, mess, args)
	default:
		l.error(ctx, fmt.Sprintf("error message %v has unknown type %v", message, mess), args)
	}
}

func (l *Logger) debug(ctx context.Context, message string, args ...interface{}) {
	if args == nil {
		l.logger.Info().Str("Request id", getRequestID(ctx)).Msg(message)
	} else {
		l.logger.Info().Str("Request id", getRequestID(ctx)).Msgf(message, args...)
	}
}

func (l *Logger) Info(ctx context.Context, message string, args ...interface{}) {
	if args == nil {
		l.logger.Info().Str("Request id", getRequestID(ctx)).Msg(message)
	} else {
		l.logger.Info().Str("Request id", getRequestID(ctx)).Msgf(message, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, message string, args ...interface{}) {
	if args == nil {
		l.logger.Warn().Str("Request id", getRequestID(ctx)).Msg(message)
	} else {
		l.logger.Warn().Str("Request id", getRequestID(ctx)).Msgf(message, args...)
	}
}

func (l *Logger) error(ctx context.Context, message string, args ...interface{}) {
	if args == nil {
		l.logger.Error().Str("Request id", getRequestID(ctx)).Msg(message)
	} else {
		l.logger.Error().Str("Request id", getRequestID(ctx)).Msgf(message, args...)
	}
}

func getRequestID(ctx context.Context) (reqID string) {
	if ctx == nil {
		return "-"
	}
	reqID, ok := ctx.Value(RequestID).(string)
	if !ok {
		return "-"
	}
	return
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
