// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/robberphex/grpc-in-memory/helloworld"
)

func main() {
	svc := NewServerImpl()
	c := serverToClient(svc)

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("world_unary_%d", i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}

	for i := 0; i < 5; i++ {
		sClient, err := c.SayHelloStream(ctx)
		if err != nil {
			log.Fatalf("could not SayHelloStream: %v", err)
		}
		defer sClient.CloseSend()
		err = sClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("world_stream_%d", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		time.Sleep(time.Millisecond * 1000)
		reply, err := sClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}
}
