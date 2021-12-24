# Protobuf Example


## Generate Go Files

```shell
$ cd protobuf
~/learn-by-test/protobuf

$ protoc -I ./proto \
    --go_out ./proto \
    --go_opt module=github.com/gojustforfun/learn-by-test/grpc/proto \
    ./proto/trip.proto
```