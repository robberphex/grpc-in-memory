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

	// unary
	for i := 0; i < 5; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("world_unary_%d", i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}

	// RequestStream
	for i := 0; i < 5; i++ {
		sClient, err := c.SayHelloRequestStream(ctx)
		if err != nil {
			log.Fatalf("could not SayHelloRequestStream: %v", err)
		}
		defer sClient.CloseSend()
		err = sClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("SayHelloRequestStream_%d", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		err = sClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("SayHelloRequestStream_%d_dup", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		time.Sleep(time.Millisecond * 1000)
		reply, err := sClient.CloseAndRecv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}

	// ReplyStream
	for i := 0; i < 5; i++ {
		sClient, err := c.SayHelloReplyStream(ctx, &pb.HelloRequest{Name: fmt.Sprintf("SayHelloReplyStream_%d", i)})
		if err != nil {
			log.Fatalf("could not SayHelloReplyStream: %v", err)
		}
		reply, err := sClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
		reply, err = sClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}

	// BiStream
	for i := 0; i < 5; i++ {
		sClient, err := c.SayHelloBiStream(ctx)
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
