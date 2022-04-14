package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
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
	// 限流入口
	e, d := sentinel.Entry("warm_Reject", sentinel.WithTrafficType(base.Inbound))
	if d != nil {
		return nil, errors.New("too many request")
	}
	logger.Infof("Received: %v", in.GetName())
	e.Exit()
	// 限流结束
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
