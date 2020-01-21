package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"google.golang.org/grpc"

	"github.com/yakumioto/example-go/grpc/hello_world/protos"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	sayHello(conn)

	sayHelloClientStream(conn)

	sayHelloServerStream(conn)

	sayHelloBidiStream(conn)
}

func sayHello(conn *grpc.ClientConn) {
	client := protos.NewGreeterClient(conn)

	reply, err := client.SayHello(context.Background(), &protos.HelloRequest{Name: "Mioto"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", reply.Message)
}

func sayHelloClientStream(conn *grpc.ClientConn) {
	client := protos.NewGreeterClient(conn)
	stream, err := client.SayHelloClientStream(context.Background())
	if err != nil {
		log.Printf("failed to get stream: %v", err)
		return
	}

	for i := 0; i < 100; i++ {
		err = stream.Send(&protos.HelloRequest{Name: "Mioto" + strconv.Itoa(i)})
		if err != nil {
			log.Printf("failed to send: %v", err)
			return
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("failed to recv: %v", err)
	}

	log.Printf("Greeting: %s", reply.Message)
}

func sayHelloServerStream(conn *grpc.ClientConn) {
	client := protos.NewGreeterClient(conn)
	replyStream, err := client.SayHelloServerStream(context.Background(), &protos.HelloRequest{Name: "Mioto"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		reply, err := replyStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("failed to recv: %v", err)
			return
		}

		log.Printf("Greeting: %s", reply.Message)
	}
}

func sayHelloBidiStream(conn *grpc.ClientConn) {
	client := protos.NewGreeterClient(conn)

	stream, err := client.SayHelloBidiStream(context.Background())
	if err != nil {
		log.Printf("failed to get stream: %v", err)
		return
	}

	for i := 0; i < 100; i++ {
		err = stream.Send(&protos.HelloRequest{Name: "Mioto" + strconv.Itoa(i)})
		if err != nil {
			log.Printf("failed to send: %v", err)
			return
		}

		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("failed to recv: %v", err)
			return
		}

		log.Printf("Greeting: %s", reply.Message)
	}
}
