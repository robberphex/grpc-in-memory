// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"

	"github.com/robberphex/grpc-in-memory/helloworld"
	pb "github.com/robberphex/grpc-in-memory/helloworld"
)

// helloworld.GreeterServer 的实现
type server struct {
	// 为了后面代码兼容，必须聚合UnimplementedGreeterServer
	// 这样以后在proto文件中新增加一个方法的时候，这段代码至少不会报错
	pb.UnimplementedGreeterServer
}

// unary调用的服务端代码
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// 客户端流式调用的服务端代码
// 接收两个req，然后返回一个resp
func (s *server) SayHelloRequestStream(streamServer pb.Greeter_SayHelloRequestStreamServer) error {
	req, err := streamServer.Recv()
	if err != nil {
		log.Printf("error receiving: %v", err)
		return err
	}
	log.Printf("Received: %v", req.GetName())
	req, err = streamServer.Recv()
	if err != nil {
		log.Printf("error receiving: %v", err)
		return err
	}
	log.Printf("Received: %v", req.GetName())
	streamServer.SendAndClose(&pb.HelloReply{Message: "Hello " + req.GetName()})
	return nil
}

// 服务端流式调用的服务端代码
// 接收一个req，然后发送两个resp
func (s *server) SayHelloReplyStream(req *pb.HelloRequest, streamServer pb.Greeter_SayHelloReplyStreamServer) error {
	log.Printf("Received: %v", req.GetName())
	err := streamServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName()})
	if err != nil {
		log.Printf("error Send: %+v", err)
		return err
	}
	err = streamServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName() + "_dup"})
	if err != nil {
		log.Printf("error Send: %+v", err)
		return err
	}
	return nil
}

// 双向流式调用的服务端代码
func (s *server) SayHelloBiStream(streamServer helloworld.Greeter_SayHelloBiStreamServer) error {
	req, err := streamServer.Recv()
	if err != nil {
		log.Printf("error receiving: %+v", err)
		// 及时将错误返回给客户端，下同
		return err
	}
	log.Printf("Received: %v", req.GetName())
	err = streamServer.Send(&pb.HelloReply{Message: "Hello " + req.GetName()})
	if err != nil {
		log.Printf("error Send: %+v", err)
		return err
	}
	// 离开这个函数后，streamServer会关闭，所以不推荐在单独的goroute发送消息
	return nil
}

// 新建一个服务端实现
func NewServerImpl() *server {
	return &server{}
}
