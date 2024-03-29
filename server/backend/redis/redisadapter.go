package redis

import (
	repository "github.com/dkrizic/todo/server/backend/repository"
	redis "github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel"
	"golang.org/x/net/context"
)

type RedisAdapter struct {
	redis *redis.Client
}

func NewRedisAdapter(redis *redis.Client) *RedisAdapter {
	return &RedisAdapter{
		redis: redis,
	}
}

func (ra *RedisAdapter) ReadFromRedis(ctx context.Context, key string) (*repository.Todo, error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "ReadFromRedis")
	defer span.End()
	res, err := ra.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if res == 0 {
		return nil, nil
	}
	data := &repository.Todo{
		Id:          key,
		Title:       ra.redis.HGet(ctx, key, title).Val(),
		Description: ra.redis.HGet(ctx, key, description).Val(),
	}
	return data, nil
}

func (ra *RedisAdapter) WriteToRedis(ctx context.Context, todo *repository.Todo) (before *repository.Todo, current *repository.Todo, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "WriteToRedis")
	defer span.End()
	before, err = ra.ReadFromRedis(ctx, todo.Id)
	if err != nil {
		return nil, nil, err
	}
	ra.redis.HSet(ctx, todo.Id, title, todo.Title)
	ra.redis.HSet(ctx, todo.Id, description, todo.Description)
	current = todo
	if err != nil {
		return nil, nil, err
	}
	return before, current, nil
}

func (ra *RedisAdapter) DeleteFromRedis(ctx context.Context, key string) (before *repository.Todo, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "DeleteFromRedis")
	defer span.End()
	before, err = ra.ReadFromRedis(ctx, key)
	if err != nil {
		return nil, err
	}
	ra.redis.Del(ctx, key)
	return before, nil
}
