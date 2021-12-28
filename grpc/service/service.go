package service

import (
	"context"

	"github.com/gojustforfun/learn-by-test/grpc/proto/gen/go/trippb"
)

/*
type TripServiceServer interface {
	GetTrip(context.Context, *GetTripRequest) (*GetTripResponse, error)
	mustEmbedUnimplementedTripServiceServer()
}
*/

type TripService struct {
	trippb.UnimplementedTripServiceServer
}

func (t *TripService) GetTrip(ctx context.Context, in *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{Id: in.Id, Trip: &trippb.Trip{
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
	}}, nil
}
