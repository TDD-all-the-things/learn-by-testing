# gRPC-Gateway

## Installation

### Protocol Buffers Compiler and Plugin

```shell
# Download compiler
brew insall protobuf
https://github.com/protocolbuffers/protobuf/releases/latest

# install compiler plugin
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

```shell
protoc
protoc-gen-go
```


### gRPC Plugin

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```shell
protoc-gen-go-grpc
```

### gRPC-Gateway Plugin

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

```shell
protoc-gen-grpc-gateway
protoc-gen-openapiv2
```

# Helloworld Example
## Generate File

### For gRPC

```shell
protoc -I ./proto \
   --go_out ./proto --go_opt paths=source_relative \
   --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
   ./proto/helloworld/hello_world.proto
```

```shell
    go run cmd/helloworld/grpc/server/main.go
    go run cmd/helloworld/grpc/client/main.go
```

### For gRPC-Gateway

1. with annotations .proto

```shell
protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
  ./proto/helloworld/hello_world.proto
```

2. without change `.proto` file, using api configuration file

```shell
protoc -I ./proto \
  --grpc-gateway_out ./proto \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt grpc_api_configuration=./proto/helloworld/hello_world.yaml \
  ./proto/helloworld/hello_world.proto
```

## Running

```shell
    go run cmd/helloworld/gateway/main.go
    curl -X POST -k http://localhost:8090/v1/example/hello -d '{"name":"Go"}'
```


# Reference

- https://github.com/grpc-ecosystem/grpc-gateway
