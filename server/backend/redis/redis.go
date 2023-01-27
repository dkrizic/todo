package redis

import (
	"context"
	repository "github.com/dkrizic/todo/server/backend/repository"
	redis "github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

const (
	title       = "title"
	description = "description"
)

type server struct {
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
	var _ repository.TodoRepository = myServer
	return myServer
}

func (*server) Name() string {
	return "redis"
}

func (s *server) Create(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Create")
	defer span.End()
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
	return &repository.CreateOrUpdateResponse{
		Todo: current,
	}, nil
}

func (s *server) Update(ctx context.Context, req *repository.CreateOrUpdateRequest) (resp *repository.CreateOrUpdateResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Create")
	defer span.End()
	llog := log.WithFields(log.Fields{
		"id":          req.Todo.Id,
		"title":       req.Todo.Title,
		"description": req.Todo.Description,
	})
	llog.Info("Updating todo")
	_, current, err := s.RedisAdapter.WriteToRedis(ctx, req.Todo)
	return &repository.CreateOrUpdateResponse{
		Todo: current,
	}, nil
}

func (s *server) GetAll(ctx context.Context, req *repository.GetAllRequest) (resp *repository.GetAllResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "GetAll")
	defer span.End()
	log.WithField("implementation", s.Name()).Info("Getting all todos")

	var cursor uint64 = 0
	todos := make([]*repository.Todo, 0)
	for {
		var keys []string
		var err error
		ctx2, span2 := otel.Tracer("redis").Start(ctx, "Scan")
		keys, cursor, err = s.RedisAdapter.redis.Scan(ctx2, cursor, "", 10).Result()
		span2.End()
		if err != nil {
			log.WithError(err).Fatal("Failed to get keys")
			span.RecordError(err)
			return nil, err
		}
		for _, key := range keys {
			log.WithField("key", key).Info("Found key")
			var title string
			var description string
			{
				ctx2, span := otel.Tracer("redis").Start(ctx, "GetAll/ReadTitle")
				span.SetAttributes(attribute.String("key", key))
				title = s.RedisAdapter.redis.HGet(ctx2, key, title).Val()
				span.End()
			}
			{
				ctx2, span := otel.Tracer("redis").Start(ctx, "GetAll/ReadDescription")
				span.SetAttributes(attribute.String("key", key))
				description = s.RedisAdapter.redis.HGet(ctx2, key, description).Val()
				span.End()
			}
			todos = append(todos, &repository.Todo{
				Id:          key,
				Title:       title,
				Description: description,
			})
		}
		if cursor == 0 {
			break
		}
	}

	return &repository.GetAllResponse{
		Todos: todos,
	}, nil
}

func (s *server) Get(ctx context.Context, req *repository.GetRequest) (resp *repository.GetResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Get")
	defer span.End()
	llog := log.WithField("id", req.Id)
	llog.WithField("id", req.Id).Info("Getting todo")
	data, err := s.RedisAdapter.ReadFromRedis(ctx, req.Id)
	if err != nil {
		llog.WithError(err).Fatal("Failed to get todo")
		return nil, err
	}
	return &repository.GetResponse{
		Todo: data,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *repository.DeleteRequest) (resp *repository.DeleteResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Get")
	defer span.End()
	log.WithField("id", req.Id).Info("Deleting todo")
	s.RedisAdapter.redis.Del(ctx, req.Id)
	return &repository.DeleteResponse{
		Id: req.Id,
	}, nil
}
