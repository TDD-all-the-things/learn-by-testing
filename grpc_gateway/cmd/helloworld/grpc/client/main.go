package main

import (
	"log"
	"time"

	pb "github.com/gojustforfun/learn-by-test/grpc_gateway/proto/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Longyue Li"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
