package mongodb_test

import (
	"context"
	"fewoserv/internal/infrastructure/config"
	"fewoserv/pkg/mongodb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	configuration = config.Load()
)

func TestNewClient(t *testing.T) {
	client, err := mongodb.NewClient(configuration.MongoDB.MongoDBUri, configuration.MongoDB.MongoDBName)

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestClient_Connect_Disconnect(t *testing.T) {
	client, err := mongodb.NewClient(configuration.MongoDB.MongoDBUri, configuration.MongoDB.MongoDBName)

	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, client.Connect(ctx))
	require.NoError(t, client.Disconnect(ctx))
}

func TestClient_GetCollection(t *testing.T) {
	client, err := mongodb.NewClient(configuration.MongoDB.MongoDBUri, configuration.MongoDB.MongoDBName)

	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, client.Connect(ctx))

	assert.NotNil(t, client.GetCollection("test-collection"))

	require.NoError(t, client.Disconnect(ctx))
}
