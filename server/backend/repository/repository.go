package repository

import (
	"context"
)

type TodoRepository interface {
	Create(ctx context.Context, req *CreateOrUpdateRequest) (resp *CreateOrUpdateResponse, err error)
	Update(ctx context.Context, req *CreateOrUpdateRequest) (resp *CreateOrUpdateResponse, err error)
	GetAll(ctx context.Context, req *GetAllRequest) (resp *GetAllResponse, err error)
	Get(ctx context.Context, req *GetRequest) (resp *GetResponse, err error)
	Delete(ctx context.Context, req *DeleteRequest) (resp *DeleteResponse, err error)
}

type Todo struct {
	Id          string
	Title       string
	Description string
	Status      string
}

type CreateOrUpdateRequest struct {
	Todo *Todo
}

type CreateOrUpdateResponse struct {
	Todo *Todo
}

type GetAllRequest struct {
	// nothing
}

type GetAllResponse struct {
	Todos []*Todo
}

type GetRequest struct {
	Id string
}

type GetResponse struct {
	Todo *Todo
}

type DeleteRequest struct {
	Id string
}

type DeleteResponse struct {
	Id string
}
