package memory

import (
	"context"
	repository "github.com/dkrizic/todo/server/backend/repository"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"golang.org/x/exp/maps"
)

var todoMap = map[string]*repository.Todo{}

type server struct {
	maxEntries int
}

func NewServer(maxEntries int) *server {
	log.WithField("maxEntries", maxEntries).Info("Creating new memory server")
	myServer := &server{
		maxEntries: maxEntries,
	}
	// ensure server implements the inteface
	var _ repository.TodoRepository = myServer
	return myServer
}

func (s *server) Create(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("memory").Start(ctx, "Create")
	defer span.End()
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Creating new todo")
	// add to map
	todoMap[req.Todo.Id] = req.Todo
	return &repository.CreateOrUpdateResponse{
		Todo: req.Todo,
	}, nil
}

func (s *server) Update(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("memory").Start(ctx, "Update")
	defer span.End()
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Updating todo")
	todoMap[req.Todo.Id] = req.Todo
	return &repository.CreateOrUpdateResponse{
		Todo: req.Todo,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *repository.GetAllRequest) (resp *repository.GetAllResponse, err error) {
	ctx, span := otel.Tracer("memory").Start(ctx, "GetAll")
	defer span.End()
	log.WithField("count", len(todoMap)).Info("Getting all todos")
	// convert map of todoMap to slice
	_, span2 := otel.Tracer("memory").Start(ctx, "GetAll/mapValues")
	todos := maps.Values(todoMap)
	defer span2.End()
	return &repository.GetAllResponse{
		Todos: todos,
	}, nil
}

func (s *server) Get(ctx context.Context, req *repository.GetRequest) (resp *repository.GetResponse, err error) {
	ctx, span := otel.Tracer("memory").Start(ctx, "Get")
	defer span.End()
	log.WithField("id", req.Id).Info("Getting todo")
	return &repository.GetResponse{
		Todo: todoMap[req.Id],
	}, nil
}

func (s *server) Delete(ctx context.Context, req *repository.DeleteRequest) (resp *repository.DeleteResponse, err error) {
	ctx, span := otel.Tracer("memory").Start(ctx, "Delete")
	defer span.End()
	log.WithField("id", req.Id).Info("Deleting todo")
	delete(todoMap, req.Id)
	return &repository.DeleteResponse{
		Id: req.Id,
	}, nil
}
