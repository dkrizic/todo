package main

import (
	"context"
	todo "github.com/dkrizic/todo/api"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log.Info("Starting app")
	cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatal("Error connecting to server")
	}
	ctx := context.Background()
	tc := todo.NewToDoServiceClient(cc)

	_, err = tc.Create(ctx, &todo.CreateOrUpdateRequest{
		Api: "v1",
		Todo: &todo.ToDo{
			Id:          uuid.New().String(),
			Title:       "Another todo",
			Description: "This is the description of the todo",
		},
	})
	if err != nil {
		log.WithError(err).Fatal("Error creating todo")
	}

	all, err := tc.GetAll(ctx, &todo.GetAllRequest{})
	if err != nil {
		log.WithError(err).Fatal("Error getting all todos")
	}
	for _, t := range all.Todos {
		log.WithField("id", t.Id).WithField("title", t.Title).Info("Got todo")
	}
}
