package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IClient interface {
		Connect(ctx context.Context) error
		Disconnect(ctx context.Context) error
		GetCollection(collection string) *mongo.Collection
	}

	// Client encapsulates a mongo client and is responsible for communication with a mongo cluster or server

	Client struct {
		client *mongo.Client
		db     string
		uri    string
	}
)

// NewClient returns a new Client
func NewClient(databaseURI, databaseName string) (IClient, error) {
	client := &Client{db: databaseName, uri: databaseURI}

	if len(client.db) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDatabaseName, "database name cannot be empty")
	}

	if len(client.uri) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDatabaseURI, "database uri cannot be empty")
	}

	return client, nil
}

// Start initializes the connection to a mongo server / cluster and returns an error is anything goes wrong
// during the process.
func (c *Client) Connect(ctx context.Context) error {
	var (
		client *mongo.Client
		err    error
	)

	opts := options.Client().ApplyURI(c.uri)

	if err = opts.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidDatabaseURI, err.Error())
	}

	if client, err = mongo.Connect(ctx, opts); err != nil {
		return fmt.Errorf("%w: %s", ErrConnection, err.Error())
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("%w: %s", ErrPing, err.Error())
	}

	c.client = client

	fmt.Println("successfully connected to DB:", c.db)

	return nil
}

// Stop terminates the mongodb connection. It returns an error if the disconnect did not work
func (c *Client) Disconnect(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetCollection returns the target collection from the selected database
func (c *Client) GetCollection(collection string) *mongo.Collection {
	return c.client.Database(c.db).Collection(collection)
}
