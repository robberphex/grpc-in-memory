// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/robberphex/grpc-in-memory/helloworld"
	"google.golang.org/grpc"
)

// 将一个服务实现转化为一个客户端
func serverToClient(svc *server) pb.GreeterClient {
	// 创建一个基于pipe的Listener
	pipe := ListenPipe()

	s := grpc.NewServer()
	// 注册Greeter服务到gRPC
	pb.RegisterGreeterServer(s, svc)
	if err := s.Serve(pipe); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// 客户端指定使用pipe作为网络连接
	clientConn, err := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 基于pipe连接，创建gRPC客户端
	c := pb.NewGreeterClient(clientConn)
	return c
}

func main() {
	svc := NewServerImpl()
	c := serverToClient(svc)

	ctx := context.Background()

	// unary调用
	for i := 0; i < 5; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("world_unary_%d", i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}

	// 客户端流式调用
	for i := 0; i < 5; i++ {
		streamClient, err := c.SayHelloRequestStream(ctx)
		if err != nil {
			log.Fatalf("could not SayHelloRequestStream: %v", err)
		}
		err = streamClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("SayHelloRequestStream_%d", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		err = streamClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("SayHelloRequestStream_%d_dup", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		reply, err := streamClient.CloseAndRecv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}

	// 服务端流式调用
	for i := 0; i < 5; i++ {
		streamClient, err := c.SayHelloReplyStream(ctx, &pb.HelloRequest{Name: fmt.Sprintf("SayHelloReplyStream_%d", i)})
		if err != nil {
			log.Fatalf("could not SayHelloReplyStream: %v", err)
		}
		reply, err := streamClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
		reply, err = streamClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}

	// 双向流式调用
	for i := 0; i < 5; i++ {
		streamClient, err := c.SayHelloBiStream(ctx)
		if err != nil {
			log.Fatalf("could not SayHelloStream: %v", err)
		}
		err = streamClient.Send(&pb.HelloRequest{Name: fmt.Sprintf("world_stream_%d", i)})
		if err != nil {
			log.Fatalf("could not Send: %v", err)
		}
		reply, err := streamClient.Recv()
		if err != nil {
			log.Fatalf("could not Recv: %v", err)
		}
		log.Println(reply.GetMessage())
	}
}
