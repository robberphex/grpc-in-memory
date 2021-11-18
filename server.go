// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/robberphex/grpc-in-memory/helloworld"
	pb "github.com/robberphex/grpc-in-memory/helloworld"

	"google.golang.org/grpc"
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

func (s *server) SayHelloRequestStream(sServer pb.Greeter_SayHelloRequestStreamServer) error {
	req, err := sServer.Recv()
	if err != nil {
		log.Fatalf("error receiving: %v", err)
	}
	log.Printf("Received: %v", req.GetName())
	req, err = sServer.Recv()
	if err != nil {
		log.Fatalf("error receiving: %v", err)
	}
	log.Printf("Received: %v", req.GetName())
	sServer.SendAndClose(&pb.HelloReply{Message: "Hello " + req.GetName()})
	return nil
}

func (s *server) SayHelloReplyStream(req *pb.HelloRequest, sServer pb.Greeter_SayHelloReplyStreamServer) error {
	log.Printf("Received: %v", req.GetName())
	err := sServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName()})
	if err != nil {
		log.Fatalf("error Send: %v", err)
	}
	err = sServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName() + "_dup"})
	if err != nil {
		log.Fatalf("error Send: %v", err)
	}
	return nil
}

// SayHelloStream implements helloworld.GreeterServer
func (s *server) SayHelloBiStream(sServer helloworld.Greeter_SayHelloBiStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		req, err := sServer.Recv()
		if err != nil {
			log.Fatalf("error receiving: %v", err)
		}
		log.Printf("Received: %v", req.GetName())
		err = sServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName()})
		if err != nil {
			log.Fatalf("error Send: %+v", err)
		}
	}()
	// you cannot leave SayHelloStream with a goroutine that try to Recv and Send
	wg.Wait()
	//return errors.New("xxx")
	return nil
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
