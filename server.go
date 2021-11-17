// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/robberphex/grpc-in-memory/helloworld"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func NewServerImpl() *server {
	return &server{}
}

func serverToClient(svc *server) pb.GreeterClient {
	pipe := ListenPipe()
	s := grpc.NewServer()
	go func() {
		pb.RegisterGreeterServer(s, svc)
		if err := s.Serve(pipe); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	clientConn, err := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewGreeterClient(clientConn)
	return c
}

func initServer(lis net.Listener) {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
