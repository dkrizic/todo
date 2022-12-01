package redis

import (
	"context"
	todo "github.com/dkrizic/todo/api"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type server struct {
	todo.UnimplementedToDoServiceServer
	redis *redis.Client
}

type Config struct {
	Host string
	Port int
	User string
	Pass string
}

func NewServer(config *Config) *server {
	llog := log.WithFields(
		log.Fields{
			"host": config.Host,
			"port": config.Port,
			"user": config.User,
		},
	)
	llog.Info("Creating new redis server")

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Pass, // no password set
		DB:       0,           // use default DB
	})
	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		llog.WithError(status.Err()).Fatal("Failed to connect to redis")
	}
	llog.Info("Connected to redis")

	myServer := &server{
		redis: rdb,
	}
	// ensure server implements the interface
	var _ todo.ToDoServiceServer = myServer
	return &server{}
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Creating new todo")
	// add to map
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: req.Todo,
	}, nil
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Updating todo")
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
		Todos: []*todo.ToDo{},
	}, nil
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	log.WithField("id", req.Id).Info("Getting todo")
	return &todo.GetResponse{
		Api:  "v1",
		Todo: &todo.ToDo{},
	}, nil
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	log.WithField("id", req.Id).Info("Deleting todo")
	return &todo.DeleteResponse{
		Api: "v1",
		Id:  req.Id,
	}, nil
}
