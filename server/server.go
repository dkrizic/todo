package main

import (
	todo "github.com/dkrizic/proto-demo/api"
	"github.com/dkrizic/proto-demo/server/memory"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer()
	todo.RegisterToDoServiceServer(s, memory.NewServer())
	reflection.Register(s)
	log.Println("Serving gRPC on :8080")
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:8080",
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

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))))

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("Service HTTP on :8090")
	log.Fatal(gwServer.ListenAndServe())
}
