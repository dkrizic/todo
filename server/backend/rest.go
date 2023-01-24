package backend

import (
	"context"
	"encoding/json"
	"github.com/dkrizic/todo/server/backend/repository"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"io/ioutil"
	"net/http"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("backend").Start(r.Context(), "todos")
	defer span.End()
	switch r.Method {
	case "GET":
		response, err := backend.Implementation.GetAll(ctx, &repository.GetAllRequest{})
		if err != nil {
			log.WithError(err).Error("Error while getting all todos")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data, err := convertTodoStructToJson(ctx, response.Todos)
		if err != nil {
			log.WithError(err).Error("Error while converting todos to json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	case "POST":
		data, err := extracaDataFromRequest(ctx, r)
		if err != nil {
			log.WithError(err).Error("Error while extracting data from request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := convertJsonToTodoStruct(ctx, data)
		if err != nil {
			log.WithError(err).Error("Error while converting json to todo struct")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response, err := backend.Implementation.Create(ctx, &repository.CreateOrUpdateRequest{
			&todo,
		})
		if err != nil {
			log.WithError(err).Error("Error while creating todo")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("backend").Start(r.Context(), "todos/{id}")
	defer span.End()
	// get id from path
	id := chi.URLParam(r, "id")
	switch r.Method {
	case "GET":
		backend.Implementation.Get(ctx, &repository.GetRequest{
			Id: id,
		})
	case "PUT":
		backend.Implementation.Update(ctx, &repository.CreateOrUpdateRequest{})
	case "DELETE":
		backend.Implementation.Delete(ctx, &repository.DeleteRequest{})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// convert request data in json format to Todo struct
func convertJsonToTodoStruct(ctx context.Context, jsonData []byte) (todo repository.Todo, err error) {
	_, span := otel.Tracer("backend").Start(ctx, "convertJsonToTodoStruct")
	defer span.End()
	err = json.Unmarshal(jsonData, &todo)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

func extracaDataFromRequest(ctx context.Context, r *http.Request) (data []byte, err error) {
	_, span := otel.Tracer("backend").Start(ctx, "extracaDataFromRequest")
	defer span.End()
	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return data, err
	}
	return data, nil
}

func convertTodoStructToJson(ctx context.Context, todo []*repository.Todo) (jsonData []byte, err error) {
	_, span := otel.Tracer("backend").Start(ctx, "convertTodoStructToJson")
	defer span.End()
	jsonData, err = json.Marshal(todo)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}