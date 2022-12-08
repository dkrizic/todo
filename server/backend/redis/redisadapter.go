package redis

import (
	"github.com/dkrizic/todo/api/todo"
	redis "github.com/go-redis/redis/v9"
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

func (ra *RedisAdapter) ReadFromRedis(ctx context.Context, key string) (*todo.ToDo, error) {
	res, err := ra.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if res == 0 {
		return nil, nil
	}
	data := &todo.ToDo{
		Id:          key,
		Title:       ra.redis.HGet(ctx, key, title).Val(),
		Description: ra.redis.HGet(ctx, key, description).Val(),
	}
	return data, nil
}

func (ra *RedisAdapter) WriteToRedis(ctx context.Context, todo *todo.ToDo) (before *todo.ToDo, current *todo.ToDo, err error) {
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

func (ra *RedisAdapter) DeleteFromRedis(ctx context.Context, key string) (before *todo.ToDo, err error) {
	before, err = ra.ReadFromRedis(ctx, key)
	if err != nil {
		return nil, err
	}
	ra.redis.Del(ctx, key)
	return before, nil
}
