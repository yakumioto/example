package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/yakumioto/example-go/grpc/hello_world/protos"
)

type server struct {}

func (s *server) SayHello(ctx context.Context, in *protos.HelloRequest) (*protos.HelloReply, error){
	return &protos.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protos.RegisterGreeterServer(s, new(server))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
