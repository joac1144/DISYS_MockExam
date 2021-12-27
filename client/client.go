package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	api "DISYS_MockExam/api"
)

type Connection struct {
	clientConn *grpc.ClientConn
	client     api.IncrementServiceClient
	context    context.Context
	port       string
}

var connections []Connection

var serverPorts = []string{":9000", ":9001", ":9002"}

func main() {
	for i := range serverPorts {
		ctx, conn, c := setupConnection(i)

		newConn := Connection{
			context:    ctx,
			clientConn: conn,
			client:     c,
			port:       serverPorts[i],
		}

		connections = append(connections, newConn)

		defer newConn.clientConn.Close()
	}

	for {
		wait(2, 5)
		increment()
	}
}

func increment() {
	printed := false

	for i, v := range connections {
		incrementResponse, err := v.client.Increment(v.context, &api.IncrementMsg{})
		if err != nil {
			log.Printf("Error calling Increment on port %s: %s", v.port, err)
			connections[i] = connections[len(connections)-1]
			connections = connections[:len(connections)-1]
		} else {
			if !printed {
				log.Println(incrementResponse.Value)
				printed = true
			}
		}
	}
}

func wait(min, max int) {
	rand.Seed(time.Now().UnixNano())
	randomDelay := rand.Intn(max-min) + min
	time.Sleep(time.Second * time.Duration(randomDelay))
}

func setupConnection(index int) (context.Context, *grpc.ClientConn, api.IncrementServiceClient) {
	conn, err := grpc.Dial(serverPorts[index], grpc.WithInsecure())

	if err != nil {
		log.Printf("Error: %v", err)
	}

	context := context.Background()

	client := api.NewIncrementServiceClient(conn)

	return context, conn, client
}
