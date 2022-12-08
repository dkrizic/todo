package redis

import (
	"context"
	"github.com/dkrizic/todo/api/todo"
	redis "github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	title       = "title"
	description = "description"
)

type server struct {
	todo.UnimplementedToDoServiceServer
	RedisAdapter *RedisAdapter
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

	redisAdapter := &RedisAdapter{
		redis: rdb,
	}

	myServer := &server{
		RedisAdapter: redisAdapter,
	}
	// ensure server implements the interface
	var _ todo.ToDoServiceServer = myServer
	return myServer
}

func (s *server) Create(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	llog := log.WithFields(log.Fields{
		"id":          req.Todo.Id,
		"title":       req.Todo.Title,
		"description": req.Todo.Description,
	})
	llog.Info("Creating todo")
	_, current, err := s.RedisAdapter.WriteToRedis(ctx, req.Todo)
	if err != nil {
		llog.WithError(err).Fatal("Failed to create todo")
		return nil, err
	}
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: current,
	}, nil
}

func (s *server) Update(ctx context.Context, req *todo.CreateOrUpdateRequest) (resp *todo.CreateOrUpdateResponse, err error) {
	llog := log.WithFields(log.Fields{
		"id":          req.Todo.Id,
		"title":       req.Todo.Title,
		"description": req.Todo.Description,
	})
	llog.Info("Updating todo")
	_, current, err := s.RedisAdapter.WriteToRedis(ctx, req.Todo)
	return &todo.CreateOrUpdateResponse{
		Api:  "v1",
		Todo: current,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *todo.GetAllRequest) (resp *todo.GetAllResponse, err error) {
	log.Info("Getting all todos")

	var cursor uint64 = 0
	todos := make([]*todo.ToDo, 0)
	for {
		var keys []string
		var err error
		keys, cursor, err = s.RedisAdapter.redis.Scan(ctx, cursor, "", 10).Result()
		if err != nil {
			log.WithError(err).Fatal("Failed to get keys")
			return nil, err
		}
		for _, key := range keys {
			log.WithField("key", key).Info("Found key")
			todos = append(todos, &todo.ToDo{
				Id:          key,
				Title:       s.RedisAdapter.redis.HGet(ctx, key, title).Val(),
				Description: s.RedisAdapter.redis.HGet(ctx, key, description).Val(),
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
	llog := log.WithField("id", req.Id)
	llog.WithField("id", req.Id).Info("Getting todo")
	data, err := s.RedisAdapter.ReadFromRedis(ctx, req.Id)
	if err != nil {
		llog.WithError(err).Fatal("Failed to get todo")
		return nil, err
	}
	return &todo.GetResponse{
		Api:  "v1",
		Todo: data,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	log.WithField("id", req.Id).Info("Deleting todo")
	s.RedisAdapter.redis.Del(ctx, req.Id)
	return &todo.DeleteResponse{
		Api: "v1",
		Id:  req.Id,
	}, nil
}
