package main

import (
	"context"
	"log"
	"net"
	"net/http"

	helloworldpb "github.com/gojustforfun/learn-by-test/grpc_gateway/proto/helloworld"
	"github.com/gojustforfun/learn-by-test/grpc_gateway/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(grpcServer, server.NewHelloworld())

	go func() {
		log.Fatalln(grpcServer.Serve(ln))
	}()

	conn, err := grpc.DialContext(context.Background(), ln.Addr().String(), grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln("Failed to dial grpc server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    "localhost:8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://" + gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())

}
