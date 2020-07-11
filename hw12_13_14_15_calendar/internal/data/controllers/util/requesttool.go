package util

import (
	"context"
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
)

const RequestID = contextKey("RequestID")
const UserID = contextKey("UserID")
const AuthHeaderKey = "x-user-id"

type contextKey string

func GetUserID(ctx context.Context) (userID string) {
	if ctx == nil {
		return
	}
	userID, _ = ctx.Value(UserID).(string)
	return
}

func SetUserID(ctx context.Context, userID string) context.Context {
	if len(GetUserID(ctx)) == 0 {
		return context.WithValue(ctx, UserID, userID)
	}
	return ctx
}

func GetRequestID(ctx context.Context) (reqID string) {
	if ctx == nil {
		return
	}
	reqID, _ = ctx.Value(RequestID).(string)
	return
}

func SetRequestID(ctx context.Context) context.Context {
	if len(GetRequestID(ctx)) == 0 {
		reqid := uuid.NewV4()
		return context.WithValue(ctx, RequestID, reqid.String())
	}
	return ctx
}

type HTTPReqInfo struct {
	IP, Start, Method, Path, Httpver, Code, Useragent string
	Latency                                           time.Duration
}

func NewHTTPReqInfo(addr string, start time.Time, method, path, httpver string, useragent string) *HTTPReqInfo {
	return &HTTPReqInfo{
		IP:        ipFromHostPort(addr),
		Start:     start.Format(entities.LayoutISO),
		Method:    method,
		Path:      path,
		Httpver:   httpver,
		Useragent: useragent,
	}
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
