package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/gojustforfun/learn-by-test/grpc/proto/gen/go/trippb"
	"github.com/gojustforfun/learn-by-test/grpc/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

func TestTripService_GetTrip(t *testing.T) {

	ln := bufconn.Listen(1024 * 1024)

	go func() {
		grpcServer := grpc.NewServer()
		trippb.RegisterTripServiceServer(grpcServer, &service.TripService{})
		assert.NoError(t, grpcServer.Serve(ln))
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return ln.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithBlock(), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := trippb.NewTripServiceClient(conn)

	id := "Go Trip"
	trip := trippb.Trip{
		Start: "XXXX",
		StartPosition: &trippb.Location{
			Latitude:  91.92,
			Longitude: 83.49,
		},
		PathPositions: []*trippb.Location{
			{Latitude: 33.21, Longitude: 44.18},
			{Latitude: 29.78, Longitude: 57.46},
		},
		End: "YYYY",
		EndPosition: &trippb.Location{
			Latitude:  70.71,
			Longitude: 63.27,
		},
		DurationInSec: 3600,
		FeeInCent:     10000,
		Status:        trippb.Status_IN_PROGRESS,
	}

	resp, err := client.GetTrip(context.Background(), &trippb.GetTripRequest{Id: id})

	assert.NoError(t, err)
	assert.Equal(t, id, resp.Id)
	assert.True(t, proto.Equal(resp.Trip, &trip))
}
