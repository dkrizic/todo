package notification

import (
	"context"
	"encoding/json"
	"github.com/dkrizic/todo/server/backend/repository"
	"github.com/dkrizic/todo/server/sender"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type server struct {
	original repository.TodoRepository
	sender   *sender.Sender
	enabled  bool
}

type NotificationConfig struct {
	Original repository.TodoRepository
	Sender   *sender.Sender
	Enabled  bool
}

func NewServer(config *NotificationConfig) *server {
	myServer := &server{
		original: config.Original,
		sender:   config.Sender,
		enabled:  config.Enabled,
	}
	// ensure server implements the interface
	var _ repository.TodoRepository = myServer
	log.Info("Notification server created")
	return myServer
}

func (s *server) Create(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "Create")
	defer span.End()
	before, err3 := s.original.Get(ctx, &repository.GetRequest{Id: req.Todo.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Create(ctx, req)
	if err == nil {
		if s.enabled {
			change := repository.Change{
				Before:     before.Todo,
				After:      resp.Todo,
				ChangeType: "CREATE",
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) Update(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "Update")
	defer span.End()
	before, err3 := s.original.Get(ctx, &repository.GetRequest{Id: req.Todo.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Update(ctx, req)
	if err == nil {
		if s.enabled {
			change := repository.Change{
				Before:     before.Todo,
				After:      resp.Todo,
				ChangeType: "UPDATE",
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) GetAll(ctx context.Context, req *repository.GetAllRequest) (resp *repository.GetAllResponse, err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "GetAll")
	defer span.End()
	log.WithField("req", req).Info("notification-GetAll")
	return s.original.GetAll(ctx, req)
}
func (s *server) Get(ctx context.Context, req *repository.GetRequest) (resp *repository.GetResponse, err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "Get")
	defer span.End()
	return s.original.Get(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *repository.DeleteRequest) (resp *repository.DeleteResponse, err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "Delete")
	defer span.End()
	before, err3 := s.original.Get(ctx, &repository.GetRequest{Id: req.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Delete(ctx, req)
	if err == nil {
		if s.enabled {
			change := repository.Change{
				Before:     before.Todo,
				After:      nil,
				ChangeType: "DELETE",
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) send(ctx context.Context, change repository.Change) (err error) {
	ctx, span := otel.Tracer("notification").Start(ctx, "send")
	defer span.End()
	data, err := convert(change)
	if err != nil {
		return err
	}
	log.WithField("change", string(data)).Info("Sending notification")
	err = s.sender.SendNotification(ctx, data)
	return nil
}

func convert(change repository.Change) (data []byte, err error) {
	data, err = json.Marshal(change)
	if err != nil {
		log.WithError(err).Error("Failed to convert change to json")
		return nil, err
	}
	return data, nil
}
