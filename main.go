// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/robberphex/grpc-in-memory/helloworld"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	pipe := ListenPipe()
	go initServer(pipe)
	clientConn, err := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer clientConn.Close()
	c := pb.NewGreeterClient(clientConn)

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
