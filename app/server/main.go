package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	echo_rpc "github.com/ekky-kharismadhany/grpc-load-balancer/server-demo/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var app_id = uuid.NewString()

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

const port = "port"

type config struct {
	port string
}

type echoServer struct {
	echo_rpc.UnimplementedEchoServerServer
}

func loadConfig() config {
	config := config{}
	config.port = os.Getenv(port)
	return config
}

func (svc echoServer) CallEcho(ctx context.Context, request *echo_rpc.Echo) (*echo_rpc.Echo, error) {
	message := request.Message
	response := fmt.Sprintf("from: %s replied by: %s", message, app_id)
	return &echo_rpc.Echo{Message: response}, nil
}

func initServer() echoServer {
	return echoServer{}
}

func main() {
	config := loadConfig()
	host := fmt.Sprintf(":%s", config.port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	reflection.Register(server)
	echo_rpc.RegisterEchoServerServer(server, initServer())
	logger.Info("server started", "host", host)
	if err := server.Serve(listener); err != nil {
		panic(err.Error())
	}
}
