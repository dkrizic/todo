package notification

import (
	"context"
	"encoding/json"
	"github.com/dkrizic/todo/api/todo"
	"github.com/dkrizic/todo/server/sender"
	"github.com/opentracing/opentracing-go"
	splog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

type server struct {
	todo.UnimplementedToDoServiceServer
	original todo.ToDoServiceServer
	sender   *sender.Sender
	enabled  bool
}

type NotificationConfig struct {
	Original todo.ToDoServiceServer
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
	var _ todo.ToDoServiceServer = myServer
	log.Info("Notification server created")
	return myServer
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	before, err3 := s.original.Get(ctx, &todo.GetRequest{Id: req.Todo.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Create(ctx, req)
	if err == nil {
		if s.enabled {
			change := todo.Change{
				Before:     before.Todo,
				After:      resp.Todo,
				ChangeType: todo.ChangeType_CREATE,
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	before, err3 := s.original.Get(ctx, &todo.GetRequest{Id: req.Todo.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Update(ctx, req)
	if err == nil {
		if s.enabled {
			change := todo.Change{
				Before:     before.Todo,
				After:      resp.Todo,
				ChangeType: todo.ChangeType_UPDATE,
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) GetAll(ctx context.Context, req *todo.GetAllRequest) (resp *todo.GetAllResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "notification-GetAll")
	defer span.Finish()
	log.WithField("req", req).Info("notification-GetAll")
	return s.original.GetAll(ctx, req)
}
func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "notification-GetAll")
	defer span.Finish()
	span.LogFields(splog.String("id", req.Id))
	return s.original.Get(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	before, err3 := s.original.Get(ctx, &todo.GetRequest{Id: req.Id})
	if err3 != nil {
		log.WithError(err3).Error("Failed to get todo before deleting")
		return nil, err3
	}
	resp, err = s.original.Delete(ctx, req)
	if err == nil {
		if s.enabled {
			change := todo.Change{
				Before:     before.Todo,
				After:      nil,
				ChangeType: todo.ChangeType_DELETE,
			}
			err2 := s.send(ctx, change)
			if err2 != nil {
				log.WithError(err2).Warn("Failed to send notification")
			}
		}
	}
	return resp, err
}

func (s *server) send(ctx context.Context, change todo.Change) (err error) {
	data, err := convert(change)
	if err != nil {
		return err
	}
	log.WithField("change", string(data)).Info("Sending notification")
	err = s.sender.SendNotification(ctx, data)
	return nil
}

func convert(change todo.Change) (data []byte, err error) {
	data, err = json.Marshal(change)
	if err != nil {
		log.WithError(err).Error("Failed to convert change to json")
		return nil, err
	}
	return data, nil
}
