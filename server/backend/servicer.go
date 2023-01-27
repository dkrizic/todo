package backend

import (
	"context"
	api "github.com/dkrizic/todo/api/todo"
	"go.opentelemetry.io/otel"
)

func NewMyApiServicer() api.DefaultApiServicer {
	return &MyApiServicer{}
}

type MyApiServicer struct {
	api.DefaultApiServicer
}

func (s *MyApiServicer) CreateTodo(ctx context.Context, todo api.Todo) (response api.ImplResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Services-CreateTodo")
	defer span.End()
	return s.DefaultApiServicer.CreateTodo(ctx, todo)
}

func (s *MyApiServicer) DeleteTodo(ctx context.Context, todoId string) (response api.ImplResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Services-DeleteTodo")
	defer span.End()
	return s.DefaultApiServicer.DeleteTodo(ctx, todoId)
}

func (s *MyApiServicer) GetAllTodos(ctx context.Context) (response api.ImplResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Services-GetAllTodos")
	defer span.End()
	return s.DefaultApiServicer.GetAllTodos(ctx)
}

func (s *MyApiServicer) GetTodo(ctx context.Context, todoId string) (response api.ImplResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Services-GetTodo")
	defer span.End()
	return s.DefaultApiServicer.GetTodo(ctx, todoId)
}

func (s *MyApiServicer) UpdateTodo(ctx context.Context, todoId string, todo api.Todo) (response api.ImplResponse, err error) {
	ctx, span := otel.Tracer("redis").Start(ctx, "Services-UpdateTodo")
	defer span.End()
	return s.DefaultApiServicer.UpdateTodo(ctx, todoId, todo)
}
