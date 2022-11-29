package main

import (
	"github.com/dkrizic/proto-demo/proto/todo"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
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
	log.Println("Serving gRPC on :8080")
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial", err)
	}

	gwmux := runtime.NewServeMux()
	err = todo.RegisterToDoServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal("Failed to register gateway", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Service HTTP on :8090")
	log.Fatal(gwServer.ListenAndServe())
}
