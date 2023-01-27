package backend

import (
	api "github.com/dkrizic/todo/api/todo"
	otlp "go.opentelemetry.io/otel"
	"net/http"
)

type MyApiRouter struct {
	api.DefaultApiRouter
}

func NewMyApiRouter() api.DefaultApiRouter {
	return &MyApiRouter{}
}

func (router *MyApiRouter) CreateTodo(w http.ResponseWriter, r *http.Request) {
	_, span := otlp.Tracer("todo").Start(r.Context(), "Router-CreateTodo")
	defer span.End()
	router.DefaultApiRouter.CreateTodo(w, r)
	return
}

func (router *MyApiRouter) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	_, span := otlp.Tracer("todo").Start(r.Context(), "Router-DeleteTodo")
	defer span.End()
	router.DefaultApiRouter.DeleteTodo(w, r)
	return
}

func (router *MyApiRouter) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	_, span := otlp.Tracer("todo").Start(r.Context(), "Router-GetAllTodos")
	defer span.End()
	router.DefaultApiRouter.GetAllTodos(w, r)
	return
}

func (router *MyApiRouter) GetTodo(w http.ResponseWriter, r *http.Request) {
	_, span := otlp.Tracer("todo").Start(r.Context(), "Router-GetTodo")
	defer span.End()
	router.DefaultApiRouter.GetTodo(w, r)
	return
}

func (router *MyApiRouter) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	_, span := otlp.Tracer("todo").Start(r.Context(), "Router-UpdateTodo")
	defer span.End()
	router.DefaultApiRouter.UpdateTodo(w, r)
	return
}

func (router *MyApiRouter) MountRoutes(r *http.ServeMux) {
	r.HandleFunc("/todos", router.GetAllTodos)
	r.HandleFunc("/todos/{todoId}", router.GetTodo)
	r.HandleFunc("/todos", router.CreateTodo)
	r.HandleFunc("/todos/{todoId}", router.UpdateTodo)
	r.HandleFunc("/todos/{todoId}", router.DeleteTodo)
}
