package notification

import (
	"context"
	"github.com/dkrizic/todo/api/todo"
	"github.com/dkrizic/todo/server/sender"
	log "github.com/sirupsen/logrus"
)

type server struct {
	todo.UnimplementedToDoServiceServer
	original     todo.ToDoServiceServer
	notification sender.Sender
	enabled      bool
}

func NewServer(original todo.ToDoServiceServer, enabled bool) *server {
	myServer := &server{
		original: original,
		enabled:  enabled,
	}
	// ensure server implements the interface
	var _ todo.ToDoServiceServer = myServer
	log.WithField("enabled", enabled).Info("Notification server created")
	return myServer
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.Info("Notifcation about creation")
	return s.original.Create(ctx, req)
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	return s.original.Update(ctx, req)
}

func (s *server) GetAll(ctx context.Context, req *todo.GetAllRequest) (resp *todo.GetAllResponse, err error) {
	return s.original.GetAll(ctx, req)
}
func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	return s.original.Get(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	return s.original.Delete(ctx, req)
}
