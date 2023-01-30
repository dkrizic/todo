package backend

import (
	"context"
	"encoding/json"
	"github.com/dkrizic/todo/server/backend/repository"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"io/ioutil"
	"net/http"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("backend").Start(r.Context(), "todos")
	defer span.End()
	switch r.Method {
	case "GET":
		log.WithField("Implementation", ActiveBackend.Implementation.Name()).Info("Getting all todos")
		response, err := ActiveBackend.Implementation.GetAll(ctx, &repository.GetAllRequest{})
		if err != nil {
			log.WithError(err).Error("Error while getting all todos")
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		span.SetAttributes(attribute.KeyValue{Key: "todos", Value: attribute.Int64Value(int64(len(response.Todos)))})
		data, err := convertTodosStructToJson(ctx, response.Todos)
		if err != nil {
			log.WithError(err).Error("Error while converting todos to json")
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	case "POST":
		data, err := extractDataFromRequest(ctx, r)
		if err != nil {
			log.WithError(err).Error("Error while extracting data from request")
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := convertJsonToTodoStruct(ctx, data)
		if err != nil {
			log.WithError(err).Error("Error while converting json to todo struct")
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response, err := ActiveBackend.Implementation.Create(ctx, &repository.CreateOrUpdateRequest{
			&todo,
		})
		if err != nil {
			log.WithError(err).Error("Error while creating todo")
			span.RecordError(err)
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
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.KeyValue{Key: "id", Value: attribute.StringValue(id)})
	switch r.Method {
	case "GET":
		log.WithField("id", id).Info("Getting todo by id")
		response, err := ActiveBackend.Implementation.Get(ctx, &repository.GetRequest{
			Id: id,
		})
		if err != nil {
			log.WithError(err).Error("Error while getting todo")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data, err := convertTodoStructToJson(ctx, response.Todo)
		if err != nil {
			log.WithError(err).Error("Error while converting todo to json")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	case "PUT":
		data, err := extractDataFromRequest(ctx, r)
		if err != nil {
			log.WithError(err).Error("Error while extracting data from request")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := convertJsonToTodoStruct(ctx, data)
		if err != nil {
			log.WithError(err).Error("Error while converting json to todo struct")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check that the id in the path matches the id in the request body
		if todo.Id != id {
			log.WithField("id", id).WithField("bodyId", todo.Id).Error("Id in path does not match id in request body")
			span.SetStatus(codes.Error, "Id in path does not match id in request body")
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.WithField("id", id).Info("Updating todo by id")
		response, err := ActiveBackend.Implementation.Update(ctx, &repository.CreateOrUpdateRequest{
			&todo,
		})
		if err != nil {
			log.WithError(err).Error("Error while updating todo")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data, err = convertTodoStructToJson(ctx, response.Todo)
		if err != nil {
			log.WithError(err).Error("Error while converting todo to json")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		span.SetStatus(codes.Ok, "Todo updated successfully")
	case "DELETE":
		log.WithField("id", id).Info("Deleting todo by id")
		_, err := ActiveBackend.Implementation.Delete(ctx, &repository.DeleteRequest{})
		if err != nil {
			log.WithError(err).Error("Error while deleting todo")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return todo, err
	}
	return todo, nil
}

func extractDataFromRequest(ctx context.Context, r *http.Request) (data []byte, err error) {
	_, span := otel.Tracer("backend").Start(ctx, "extracaDataFromRequest")
	defer span.End()
	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return data, err
	}
	return data, nil
}

func convertTodosStructToJson(ctx context.Context, todos []*repository.Todo) ([]byte, error) {
	_, span := otel.Tracer("backend").Start(ctx, "convertTodosStructToJson")
	defer span.End()
	jsonData, err := json.Marshal(todos)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return jsonData, nil
}

func convertTodoStructToJson(ctx context.Context, todo *repository.Todo) (jsonData []byte, err error) {
	_, span := otel.Tracer("backend").Start(ctx, "convertTodoStructToJson")
	defer span.End()
	jsonData, err = json.Marshal(todo)
	if err != nil {
		span.RecordError(err)
		return jsonData, err
	}
	return jsonData, nil
}
