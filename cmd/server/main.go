package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	init "micro_server_kit_golang/initialize"

	"google.golang.org/grpc"
	pb "micro_server_kit_golang/proto"
)

var (
	port   = init.Cfg.ServerConfig.Port
	logger = init.Logger.Sugar()
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Infof("Received: %v", in.GetName())
	// init.MysqlDb do something
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	{
		init.ConsulRegister("ip/addr", port)
		defer init.ConsulDeregister()
	}
	logger.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
