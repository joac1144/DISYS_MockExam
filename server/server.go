package main

import (
	"context"
	"log"
	"net"
	"os"

	api "DISYS_MockExam/api"

	"google.golang.org/grpc"
)

type Server struct {
	api.UnimplementedIncrementServiceServer
}

var counter int32

func main() {

	args := os.Args

	port := ":" + args[1]

	setupServer(port)
}

func (server *Server) Increment(context context.Context, in *api.IncrementMsg) (*api.Response, error) {

	counter++

	return &api.Response{Value: counter - 1}, nil
}

func setupServer(port string) {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	server := grpc.NewServer()

	api.RegisterIncrementServiceServer(server, &Server{})

	log.Printf("Listening at: %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
