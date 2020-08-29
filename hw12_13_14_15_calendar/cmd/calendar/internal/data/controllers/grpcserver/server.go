//go:generate  protoc -I. -I../../../../api grpcapi/api.proto --go_out=plugins=grpc:../../../../api/grpcapi --grpc-gateway_out=logtostderr=true:../../../../api

package grpcserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	api "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/api/grpcapi"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar/internal/domain/usecases"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/data/controllers/util"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	mainusecase "github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

var headers = []string{
	util.AuthHeaderKey,
}

type GRPCServer struct {
	logger   mainusecase.Logger
	wg       *sync.WaitGroup
	calendar usecases.Calendar
	server   *grpc.Server
	gwserver *http.Server
}

func NewGRPCServer(wg *sync.WaitGroup, logger mainusecase.Logger, calendar usecases.Calendar) *GRPCServer {
	return &GRPCServer{
		logger:   logger,
		wg:       wg,
		calendar: calendar}
}

func (cs *GRPCServer) AddEvent(ctx context.Context, req *api.AddEventRequest) (*api.AddEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	id, err := cs.calendar.MakeEvent(nctx, req.GetTitle(), req.GetDatetime(), req.GetText(), userid, req.GetDuration(), req.GetTimenotify())
	if err != nil {
		if errors.Is(err, entities.ErrDateBusy) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Aborted, err.Error())
	}

	resp := &api.AddEventResponse{
		Result: &api.AddEventResponse_Id{
			Id: id,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	id, err := cs.calendar.DeleteEvent(nctx, userid, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	resp := &api.DeleteEventResponse{
		Result: &api.DeleteEventResponse_Id{
			Id: id,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	id, err := cs.calendar.UpdateEvent(nctx, userid, req.GetId(), req.GetTitle(), req.GetDatetime(), req.GetText(), req.GetDuration(), req.GetTimenotify())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	resp := &api.UpdateEventResponse{
		Result: &api.UpdateEventResponse_Id{
			Id: id,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) GetDateEvent(ctx context.Context, req *api.GetDateEventRequest) (*api.GetDateEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	events, err := cs.calendar.GetDateEvents(nctx, req.GetDate(), userid)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	evs, err := toPBEvents(events)
	if err != nil {
		return nil, err
	}

	resp := &api.GetDateEventResponse{
		Result: &api.GetDateEventResponse_Events{
			Events: evs,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) GetWeekEvent(ctx context.Context, req *api.GetWeekEventRequest) (*api.GetWeekEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	events, err := cs.calendar.GetWeekEvents(nctx, req.GetDate(), userid)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	evs, err := toPBEvents(events)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.GetWeekEventResponse{
		Result: &api.GetWeekEventResponse_Events{
			Events: evs,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) GetMonthEvent(ctx context.Context, req *api.GetMonthEventRequest) (*api.GetMonthEventResponse, error) {
	nctx := util.SetRequestID(ctx)
	userid := util.GetUserID(ctx)

	events, err := cs.calendar.GetMonthEvents(nctx, req.GetDate(), userid)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	evs, err := toPBEvents(events)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.GetMonthEventResponse{
		Result: &api.GetMonthEventResponse_Events{
			Events: evs,
		},
	}
	return resp, nil
}

func (cs *GRPCServer) ServeGW(addr string, addrgw string) {
	defer cs.wg.Done()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cs.logger.Info(ctx, "starting grpc gateway server at %v", addrgw)

	mux := runtime.NewServeMux(
		runtime.WithMetadata(injectHeadersIntoMetadata),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := api.RegisterCalendarServiceHandlerFromEndpoint(ctx, mux, addr, opts)
	if err != nil {
		cs.logger.Error(ctx, errors.Wrapf(err, "can't register gateway from grpc endpoint at addr %v", addr))
		return
	}
	cs.gwserver = &http.Server{
		Addr:    addrgw,
		Handler: mux,
	}

	if err := cs.gwserver.ListenAndServe(); err != http.ErrServerClosed {
		cs.logger.Error(ctx, errors.Wrapf(err, "can't start  grpc gateway server at %v", addrgw))
	}
}

func (cs *GRPCServer) StopGWServe() {
	ctx := context.Background()
	cs.logger.Info(ctx, "stopping grpc gw server")
	defer cs.logger.Info(ctx, "grpc gw stopped")
	if cs.gwserver == nil {
		cs.logger.Error(ctx, "grpc gw server is nil")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := cs.gwserver.Shutdown(ctx); err != nil {
		cs.logger.Error(ctx, "can't stop grpc gw server with error: %v", err)
	}
}

func (cs *GRPCServer) PrepareGRPCListener(addr string) net.Listener {
	cs.logger.Info(context.Background(), "GRPC server: starting tcp listener at %v", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		cs.logger.Error(context.Background(), errors.Wrapf(err, "GRPC server: can't start tcp listening at addr %v", addr))
		return nil
	}
	return l
}

func (cs *GRPCServer) Serve(listener net.Listener) {
	defer cs.wg.Done()
	cs.logger.Info(context.Background(), "starting grpc server at %v", listener.Addr().String())

	cs.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(cs.loggingUnary, cs.authUnary)),
	)
	api.RegisterCalendarServiceServer(cs.server, cs)

	if err := cs.server.Serve(listener); err != http.ErrServerClosed {
		cs.logger.Error(context.Background(), errors.Wrapf(err, "can't start grpc server at %v", listener.Addr().String()))
	}
}

func (cs *GRPCServer) StopServe() {
	ctx := context.Background()
	cs.logger.Info(ctx, "stopping grpc server")
	defer cs.logger.Info(ctx, "grpc server stopped")

	cs.server.GracefulStop()
}

func (cs *GRPCServer) authUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	userID := ""

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if uID, ok := md[util.AuthHeaderKey]; ok {
			userID = strings.Join(uID, ",")
		} else {
			return nil, status.Error(codes.PermissionDenied, "unauthorized user")
		}
	}
	newCtx := util.SetUserID(ctx, userID)
	cs.logger.Info(ctx, "user id:%v incoming request", userID)
	return handler(newCtx, req)
}

func (cs *GRPCServer) loggingUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	clientIP := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}

	useragent := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua, ok := md["user-agent"]; ok {
			useragent = strings.Join(ua, ",")
		}
	}
	ri := util.NewHTTPReqInfo(clientIP, start, info.FullMethod, "", "proto3", useragent)

	newctx := util.SetRequestID(ctx)
	h, err := handler(newctx, req)
	// after executing rpc
	s, _ := status.FromError(err)
	ri.Code = s.Code().String()
	ri.Latency = time.Since(start)
	//logging
	cs.logRequest(newctx, ri)
	return h, err
}

func (cs *GRPCServer) logRequest(ctx context.Context, ri *util.HTTPReqInfo) {
	cs.logger.Info(ctx, "%s [%s] %s %s %s %s %s [%s]", ri.IP, ri.Start, ri.Method, ri.Path, ri.Httpver, ri.Code, ri.Latency, ri.Useragent)
}

func toPBEvents(events []*entities.Event) (*api.Events, error) {
	grpcEvents := make([]*api.Event, 0, len(events))
	for _, event := range events {
		pbEvent, err := toPBEvent(event)
		if err != nil {
			return nil, err
		}
		grpcEvents = append(grpcEvents, pbEvent)
	}
	return &api.Events{Event: grpcEvents}, nil
}

func toPBEvent(event *entities.Event) (*api.Event, error) {
	pbe := &api.Event{
		Id:     event.ID,
		Title:  event.Title,
		Text:   event.Text,
		Userid: event.UserID,
	}
	pbdt, err := ptypes.TimestampProto(event.DateTime)
	if err != nil {
		return nil, err
	}
	pbe.Datetime = pbdt
	pbd := ptypes.DurationProto(event.Duration)
	pbe.Duration = pbd
	pbn := ptypes.DurationProto(event.TimeNotify)
	pbe.Timenotify = pbn
	return pbe, nil
}
func injectHeadersIntoMetadata(ctx context.Context, req *http.Request) metadata.MD {
	pairs := make([]string, 0, len(headers))
	for _, h := range headers {
		if v := req.Header.Get(h); len(v) > 0 {
			pairs = append(pairs, h, v)
		}
	}
	return metadata.Pairs(pairs...)
}
