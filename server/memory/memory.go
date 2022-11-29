package memory

import (
	"context"
	todo "github.com/dkrizic/proto-demo/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

var todoMap = map[string]*todo.ToDo{}

type server struct {
	todo.UnimplementedToDoServiceServer
}

func NewServer() *server {
	log.Info("Creating new server")
	myServer := &server{}
	// ensure server implements the inteface
	var _ todo.ToDoServiceServer = myServer
	return &server{}
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Creating new todo")
	// add to map
	todoMap[req.Todo.Id] = req.Todo
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: req.Todo,
	}, nil
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Updating todo")
	todoMap[req.Todo.Id] = req.Todo
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: req.Todo,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *todo.GetAllRequest) (resp *todo.GetAllResponse, err error) {
	log.Info("Getting all todos")
	// convert map of todoMap to slice
	return &todo.GetAllResponse{
		Api:   "v1",
		Todos: maps.Values(todoMap),
	}, nil
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	log.WithField("id", req.Id).Info("Getting todo")
	return &todo.GetResponse{}, nil
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	log.WithField("id", req.Id).Info("Deleting todo")
	return &todo.DeleteResponse{}, nil
}
