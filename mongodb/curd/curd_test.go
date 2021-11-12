package curd_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCURDTestSuite(t *testing.T) {
	suite.Run(t, new(CURDTestSuite))
}

type CURDTestSuite struct {
	suite.Suite
	ctx        context.Context
	cancelFunc context.CancelFunc

	coll *mongo.Collection
}

func (s *CURDTestSuite) SetupSuite() {
	s.ctx, s.cancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27018"))
	s.NoError(err)
	s.coll = client.Database("tea").Collection("ratings")

	docs := []interface{}{
		bson.D{bson.E{Key: "type", Value: "Masala"}, bson.E{Key: "rating", Value: 10}},
		bson.D{bson.E{Key: "type", Value: "Matcha"}, bson.E{Key: "rating", Value: 7}},
		bson.D{bson.E{Key: "type", Value: "Assam"}, bson.E{Key: "rating", Value: 4}},
		bson.D{bson.E{Key: "type", Value: "Oolong"}, bson.E{Key: "rating", Value: 9}},
		bson.D{bson.E{Key: "type", Value: "Chrysanthemum"}, bson.E{Key: "rating", Value: 5}},
		bson.D{bson.E{Key: "type", Value: "Earl Grey"}, bson.E{Key: "rating", Value: 8}},
		bson.D{bson.E{Key: "type", Value: "Jasmine"}, bson.E{Key: "rating", Value: 3}},
		bson.D{bson.E{Key: "type", Value: "English Breakfast"}, bson.E{Key: "rating", Value: 6}},
		bson.D{bson.E{Key: "type", Value: "White Peony"}, bson.E{Key: "rating", Value: 4}},
	}
	_, err = s.coll.InsertMany(context.TODO(), docs)
	s.NoError(err)

}

func (s *CURDTestSuite) TearDownSuite() {
	s.coll.Drop(s.ctx)
	s.coll.Database().Client().Disconnect(s.ctx)
	s.cancelFunc()
}

func (s *CURDTestSuite) TestNothing() {
	s.Equal(true, true)
}

func (s *CURDTestSuite) TestCreate() {
	filter := bson.D{bson.E{Key: "rating", Value: bson.D{bson.E{Key: "$lt", Value: 6}}}}
	count, err := s.coll.CountDocuments(context.TODO(), filter)
	s.NoError(err)
	s.Equal(int64(4), count)
}

func (s *CURDTestSuite) TestQuery() {
	query := bson.D{bson.E{Key: "type", Value: "Assam"}}
	cursor, err := s.coll.Find(context.Background(), query)
	s.NoError(err)
	s.NotNil(cursor)
	type Res struct {
		//ID     primitive.ObjectID `bson:"_id,omitempty"`
		Type   string `bson:"type,omitempty"`
		Rating int64  `bson:"rating,omitempty"`
	}
	var res []Res
	err = cursor.All(context.TODO(), &res)
	fmt.Println(res)
	s.NoError(err)
	s.Equal([]Res{{Type: "Assam", Rating: int64(4)}}, res)

}
