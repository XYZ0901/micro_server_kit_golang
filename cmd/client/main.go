package main

import (
	"context"
	"fmt"
	"log"
	"micro_server_kit_golang/utils"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	init "micro_server_kit_golang/initialize"
	pb "micro_server_kit_golang/proto"
)

var (
	name   = init.Cfg.Name
	logger = init.Logger.Sugar()
)

func main() {
	agentSrvs, err := init.ConsulFilterServices([]string{"tag1", "tag2"}, "serverName")
	if err != nil {
		log.Fatalln(err)
	}
	as := utils.ChooseServerFromMap(agentSrvs)
	addr := fmt.Sprintf("%s:%d", as.Address, as.Port)

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}
	logger.Infof("Greeting: %s", r.GetMessage())
}
