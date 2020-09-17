package internalgrpc

import (
	"context"
	"database/sql"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

//go:generate protoc EventService.proto --go_out=plugins=grpc:. -I ../../../api

// Service grpc events service.
type Service struct {
	app app.Application
}

// New returns grpc service.
func New(a app.Application) *Service {
	return &Service{app: a}
}

// GetEvent returns event by id.
func (s *Service) GetEvent(ctx context.Context, req *GetEventRequest) (*GetEventResponse, error) {
	e, err := s.app.GetEvent(ctx, req.Id)
	switch err {
	case nil:
	case app.ErrEventIDZero:
		return nil, status.Error(codes.InvalidArgument, "wrong event id")
	case sql.ErrNoRows:
		return nil, status.Error(codes.NotFound, "event not found")
	default:
		return nil, status.Error(codes.Internal, "unable to get event")
	}
	sd, err := ptypes.TimestampProto(e.StartDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}
	ed, err := ptypes.TimestampProto(e.EndDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}

	return &GetEventResponse{
		Id:           e.ID,
		Title:        e.Title,
		StartDate:    sd,
		EndDate:      ed,
		Descr:        e.Descr,
		OwnerId:      e.OwnerID,
		NotifyBefore: e.NotifyBefore,
	}, nil
}

// GetEventList returns list of events.
func (s *Service) GetEventList(ctx context.Context, _ *empty.Empty) (*ListEventResponse, error) {
	eList, err := s.app.GetEventList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get events")
	}
	res := make([]*ListEventItem, 0, len(eList))
	for _, e := range eList {
		sd, err := ptypes.TimestampProto(e.StartDate)
		if err != nil {
			return nil, status.Error(codes.Internal, "unable to convert event")
		}
		ed, err := ptypes.TimestampProto(e.EndDate)
		if err != nil {
			return nil, status.Error(codes.Internal, "unable to convert event")
		}
		res = append(res, &ListEventItem{
			Id:           e.ID,
			Title:        e.Title,
			StartDate:    sd,
			EndDate:      ed,
			Descr:        e.Descr,
			OwnerId:      e.OwnerID,
			NotifyBefore: e.NotifyBefore,
		})
	}

	return &ListEventResponse{
		Results: res,
	}, nil
}

// CreateEvent creates new event.
func (s *Service) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	e, err := s.app.CreateEvent(ctx,
		req.Title,
		req.StartDate.AsTime(),
		req.EndDate.AsTime(),
		req.OwnerId,
		req.Descr,
		req.NotifyBefore,
	)
	switch err {
	case nil:
	case app.ErrEventFieldWrong:
		return nil, status.Error(codes.InvalidArgument, "bad request")
	default:
		return nil, status.Error(codes.Internal, "unable to create event")
	}
	sd, err := ptypes.TimestampProto(e.StartDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}
	ed, err := ptypes.TimestampProto(e.EndDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}

	return &CreateEventResponse{
		Id:           e.ID,
		Title:        e.Title,
		StartDate:    sd,
		EndDate:      ed,
		Descr:        e.Descr,
		OwnerId:      e.OwnerID,
		NotifyBefore: e.NotifyBefore,
	}, nil
}

// UpdateEvent updates event.
func (s *Service) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {
	e, err := s.app.UpdateEvent(ctx,
		req.Id,
		req.Title,
		req.StartDate.AsTime(),
		req.EndDate.AsTime(),
		req.OwnerId,
		req.Descr,
		req.NotifyBefore,
	)
	switch err {
	case nil:
	case app.ErrEventIDZero:
		return nil, status.Error(codes.InvalidArgument, "wrong event id")
	default:
		return nil, status.Error(codes.Internal, "unable to update event")
	}
	sd, err := ptypes.TimestampProto(e.StartDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}
	ed, err := ptypes.TimestampProto(e.EndDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}

	return &UpdateEventResponse{
		Id:           e.ID,
		Title:        e.Title,
		StartDate:    sd,
		EndDate:      ed,
		Descr:        e.Descr,
		OwnerId:      e.OwnerID,
		NotifyBefore: e.NotifyBefore,
	}, nil
}

// DeleteEvent deletes event.
func (s *Service) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*empty.Empty, error) {
	err := s.app.DeleteEvent(ctx, req.Id)
	switch err {
	case nil:
	case app.ErrEventIDZero:
		return nil, status.Error(codes.InvalidArgument, "wrong event id")
	default:
		return nil, status.Error(codes.Internal, "unable to delete event")
	}

	return &empty.Empty{}, nil
}
