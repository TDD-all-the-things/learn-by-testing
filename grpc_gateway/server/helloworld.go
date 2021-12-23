package server

import (
	"context"

	pb "github.com/gojustforfun/learn-by-test/grpc_gateway/proto/helloworld"
)

type helloworld struct {
	pb.UnimplementedGreeterServer
}

func NewHelloworld() *helloworld {
	return &helloworld{}
}

func (s *helloworld) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
