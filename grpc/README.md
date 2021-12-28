# gRPC Example


## Generate Go Files

### Full Path

`grpc/proto/trip.proto`

```proto
syntax = "proto3";
package trip;

option go_package = "github.com/gojustforfun/learn-by-test/grpc/proto/gen/go/trippb";
.....
```

```shell
$ cd grpc
$ pwd
~/learn-by-test/grpc

$ protoc -I ./proto \
    --go_out ./proto \
    --go_opt module=github.com/gojustforfun/learn-by-test/grpc/proto \
    --go-grpc_out ./proto \
    --go-grpc_opt module=github.com/gojustforfun/learn-by-test/grpc/proto \
    ./proto/trip.proto
```
### Relative Path

`grpc/proto/trip.proto`

```proto
syntax = "proto3";
package trip;

option go_package = "/gen/go/trippb";
.....
```

```shell
$ cd grpc
$ pwd
~/learn-by-test/grpc

$ protoc -I ./proto \
    --go_out ./proto \
    --go-grpc_out ./proto \
    ./proto/trip.proto
```