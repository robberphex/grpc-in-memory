// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/robberphex/grpc-in-memory/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	svc := NewServerImpl()
	c := serverToClient(svc)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
