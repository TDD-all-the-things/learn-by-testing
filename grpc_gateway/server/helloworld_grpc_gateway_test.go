package server_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	pb "github.com/TDD-all-the-things/learn-by-testing/grpc_gateway/proto/helloworld"
	"github.com/TDD-all-the-things/learn-by-testing/grpc_gateway/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestHelloworldGRPCGatewaySuite(t *testing.T) {
	suite.Run(t, new(HelloworldGRPCGatewaySuite))
}

type HelloworldGRPCGatewaySuite struct {
	suite.Suite
	lis            *bufconn.Listener
	gRPCServer     *grpc.Server
	gwServer       *http.Server
	grpcClientConn *grpc.ClientConn
}

func (s *HelloworldGRPCGatewaySuite) SetupSuite() {
	const bufSize = 1024 * 1024
	s.lis = bufconn.Listen(bufSize)

	s.gRPCServer = grpc.NewServer()
	pb.RegisterGreeterServer(s.gRPCServer, server.NewHelloworld())
	go func() {
		err := s.gRPCServer.Serve(s.lis)
		s.NoError(err, "Failed to start gRPC Server")
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return s.lis.Dial()
	}

	ctx := context.Background()
	var err error
	s.grpcClientConn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	s.NoError(err, "Failed to dial bufnet")

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, s.grpcClientConn)
	s.NoError(err, "Failed to register gRPC-Gateway")

	s.gwServer = &http.Server{
		Addr:    "localhost:8091",
		Handler: gwmux,
	}

	go func() {
		if err := s.gwServer.ListenAndServe(); err != http.ErrServerClosed {
			s.NoError(err, "Failed to start gRPC-Gateway proxy")
		}
	}()
}

func (s *HelloworldGRPCGatewaySuite) TestSayHello() {

	name := "gRPC"
	body := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, name))
	URL := "http://" + s.gwServer.Addr + "/v1/example/hello"

	resp, err := http.Post(URL, "application/json; charset=utf-8", body)
	s.NoError(err)
	message, err := io.ReadAll(resp.Body)
	s.NoError(err)
	s.JSONEq(fmt.Sprintf(`{"message":"Hello %s"}`, name), string(message))
}

func (s *HelloworldGRPCGatewaySuite) TearDownSuite() {
	s.gwServer.Close()
	s.grpcClientConn.Close()
	s.gRPCServer.GracefulStop()
}
