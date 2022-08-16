package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/TDD-all-the-things/learn-by-testing/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestDockerRunMongoDB(t *testing.T) {
	t.SkipNow()

	image, ports, env := "mongo:5.0.3", []string{"27017/tcp"}, []string{}
	remove, info, err := docker.Run(image, ports, env)
	require.NoError(t, err)
	defer remove()

	hostname := fmt.Sprintf("%s:%s", info.Hosts[ports[0]][0].IP, info.Hosts[ports[0]][0].Port)
	// mongodb://[uername:[password]]@localhost:27017
	// username, password := "", ""
	// uri := fmt.Sprintf("mongodb://%s:%s@%s", username, password, hostname)
	uri := fmt.Sprintf("mongodb://%s", hostname)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	defer func() {
		assert.NoError(t, client.Disconnect(context.TODO()))
	}()

	require.NoError(t, client.Ping(context.TODO(), readpref.Primary()))

	databaseNames, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"admin", "config", "local"}, databaseNames)
}
