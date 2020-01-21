package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/yakumioto/example-go/grpc/hello_world/protos"
)

type server struct{}

func (s *server) SayHello(_ context.Context, in *protos.HelloRequest) (*protos.HelloReply, error) {
	return &protos.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloClientStream(cs protos.Greeter_SayHelloClientStreamServer) error {
	names := make([]string, 0)

	for {
		in, err := cs.Recv()
		if err != nil {
			if err == io.EOF {
				return cs.SendAndClose(&protos.HelloReply{Message: "Hello " + strings.Join(names, ", ")})
			}
			log.Printf("failed to recv: %v", err)
			return err
		}

		names = append(names, in.Name)
	}
}

func (s *server) SayHelloServerStream(in *protos.HelloRequest, gss protos.Greeter_SayHelloServerStreamServer) error {
	name := in.Name

	for i := 0; i < 100; i++ {
		if err := gss.Send(&protos.HelloReply{Message: "Hello " + name + strconv.Itoa(i)}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) SayHelloBidiStream(gss protos.Greeter_SayHelloBidiStreamServer) error {
	for {
		in, err := gss.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("failed to recv: %v", err)
			break
		}

		if err := gss.Send(&protos.HelloReply{Message: "Hello " + in.Name}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protos.RegisterGreeterServer(s, new(server))

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
