package main

import (
	"github.com/dkrizic/proto-demo/proto/todo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	todo.UnimplementedToDoServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) CreateTodo(ctx context.Context, req *todo.CreateOrUpdateRequest) (*todo.CreateOrUpdateResponse, error) {
	return &todo.CreateOrUpdateResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer()
	todo.RegisterToDoServiceServer(s, NewServer())
	log.Println("Listening on :8080")
	log.Fatal(s.Serve(lis))
}
