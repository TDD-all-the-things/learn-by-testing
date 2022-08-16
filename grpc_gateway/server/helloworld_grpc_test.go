package server_test

import (
	"context"
	"net"
	"testing"

	pb "github.com/TDD-all-the-things/learn-by-testing/grpc_gateway/proto/helloworld"
	"github.com/TDD-all-the-things/learn-by-testing/grpc_gateway/server"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestHelloworldGRPCSuite(t *testing.T) {
	suite.Run(t, new(HelloworldGRPCSuite))
}

type HelloworldGRPCSuite struct {
	suite.Suite
	lis        *bufconn.Listener
	gRPCServer *grpc.Server
}

func (s *HelloworldGRPCSuite) SetupSuite() {
	const bufSize = 1024 * 1024
	s.lis = bufconn.Listen(bufSize)
	s.gRPCServer = grpc.NewServer()
	pb.RegisterGreeterServer(s.gRPCServer, server.NewHelloworld())
	go func() {
		err := s.gRPCServer.Serve(s.lis)
		s.NoError(err, "Failed to start gRPC Server")
	}()
}

func (s *HelloworldGRPCSuite) TestSayHello() {

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return s.lis.Dial()
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	s.NoError(err, "Failed to dial bufnet")
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	name := "gRPC"
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	s.NoError(err, "SayHello failed")
	s.Equal("Hello "+name, resp.GetMessage())
}

func (s *HelloworldGRPCSuite) TearDownSuite() {
	s.gRPCServer.GracefulStop()
}
