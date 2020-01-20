package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/yakumioto/example-go/grpc/hello_world/protos"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer func() {_ = conn.Close()}()

	client := protos.NewGreeterClient(conn)

	result , err := client.SayHello(context.Background(), &protos.HelloRequest{Name: "Mioto"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", result.Message)
}
