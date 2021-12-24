package trippb_test

import (
	"encoding/json"
	"testing"

	"github.com/gojustforfun/learn-by-test/protobuf/proto/gen/go/trippb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestTripSerializeAndDeserialize(t *testing.T) {

	origin := trippb.Trip{
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
		Status: trippb.Status_IN_PROGRESS,
	}

	t.Run("ProtoBuf", func(t *testing.T) {

		b, err := proto.Marshal(&origin)
		assert.NoError(t, err)

		var trip trippb.Trip
		err = proto.Unmarshal(b, &trip)
		assert.NoError(t, err)

		assert.True(t, proto.Equal(&origin, &trip))

	})

	t.Run("JSON", func(t *testing.T) {

		b, err := json.Marshal(&origin)
		assert.NoError(t, err)

		var trip trippb.Trip
		err = json.Unmarshal(b, &trip)
		assert.NoError(t, err)
		
		assert.True(t, proto.Equal(&origin, &trip))
	})

}
