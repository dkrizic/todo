package memory

import (
	"context"
	"github.com/dkrizic/todo/api/todo"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"golang.org/x/exp/maps"
)

var todoMap = map[string]*todo.ToDo{}

type server struct {
	todo.UnimplementedToDoServiceServer
	maxEntries int
}

func NewServer(maxEntries int) *server {
	log.WithField("maxEntries", maxEntries).Info("Creating new memory server")
	myServer := &server{
		maxEntries: maxEntries,
	}
	// ensure server implements the inteface
	var _ todo.ToDoServiceServer = myServer
	return myServer
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
	ctx, span := otel.Tracer("memory").Start(ctx, "GetAll")
	defer span.End()
	log.WithField("count", len(todoMap)).Info("Getting all todos")
	// convert map of todoMap to slice
	_, span2 := otel.Tracer("memory").Start(ctx, "GetAll/mapValues")
	todos := maps.Values(todoMap)
	defer span2.End()
	return &todo.GetAllResponse{
		Api:   "v1",
		Todos: todos,
	}, nil
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	log.WithField("id", req.Id).Info("Getting todo")
	return &todo.GetResponse{
		Api:  "v1",
		Todo: todoMap[req.Id],
	}, nil
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	log.WithField("id", req.Id).Info("Deleting todo")
	delete(todoMap, req.Id)
	return &todo.DeleteResponse{
		Api: "v1",
		Id:  req.Id,
	}, nil
}
