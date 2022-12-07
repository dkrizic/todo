package redis

import (
	"context"
	todo "github.com/dkrizic/todo/api"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	title       = "title"
	description = "description"
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
	return myServer
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Creating new todo")
	log.Info("Redis server: ", s.redis)
	s.redis.HSet(ctx, req.Todo.Id, title, req.Todo.Title)
	s.redis.HSet(ctx, req.Todo.Id, description, req.Todo.Description)
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: req.Todo,
	}, nil
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	log.WithField("id", req.Todo.Id).WithField("title", req.Todo.Title).Info("Updating todo")
	s.redis.HSet(ctx, req.Todo.Id, title, req.Todo.Title)
	s.redis.HSet(ctx, req.Todo.Id, description, req.Todo.Description)
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: req.Todo,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *todo.GetAllRequest) (resp *todo.GetAllResponse, err error) {
	log.Info("Getting all todos")

	var cursor uint64 = 0
	todos := make([]*todo.ToDo, 0)
	for {
		var keys []string
		var err error
		keys, cursor, err = s.redis.Scan(ctx, cursor, "", 10).Result()
		if err != nil {
			log.WithError(err).Fatal("Failed to get keys")
			return nil, err
		}
		for _, key := range keys {
			log.WithField("key", key).Info("Found key")
			todos = append(todos, &todo.ToDo{
				Id:          key,
				Title:       s.redis.HGet(ctx, key, title).Val(),
				Description: s.redis.HGet(ctx, key, description).Val(),
			})
		}
		if cursor == 0 {
			break
		}
	}

	return &todo.GetAllResponse{
		Api:   "v1",
		Todos: todos,
	}, nil
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.GetResponse, err error) {
	log.WithField("id", req.Id).Info("Getting todo")
	redis := &todo.ToDo{
		Id:          req.Id,
		Title:       s.redis.HGet(ctx, req.Id, title).Val(),
		Description: s.redis.HGet(ctx, req.Id, description).Val(),
	}
	return &todo.GetResponse{
		Api:  "v1",
		Todo: redis,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	log.WithField("id", req.Id).Info("Deleting todo")
	s.redis.Del(ctx, req.Id)
	return &todo.DeleteResponse{
		Api: "v1",
		Id:  req.Id,
	}, nil
}

func (redis *server) read(ctx context.Context, key string) (string, error) {
	return redis.redis.Get(ctx, key).Result()
}
