package mongodb

import (
	"context"
	"fewoserv/internal/infrastructure/logger"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.New("REPOSITORY")

// Repository is adapter for performing CRUD db operations
type (
	Repository[T any] struct {
		Client         IClient
		collectionName string
	}
)

func NewID() string {
	newID := primitive.NewObjectID()

	return newID.Hex()
}

// NewRepository returns new Repository with generic param
func NewRepository[T any](client IClient, collectionName string) *Repository[T] {
	return &Repository[T]{
		Client:         client,
		collectionName: collectionName,
	}
}

// DeleteByID deletes entry that matches id
func (a *Repository[_]) DeleteByID(ctx context.Context, id string) error {
	collection := a.Client.GetCollection(a.collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// Find find slice of entries by filter values
func (a *Repository[T]) Find(ctx context.Context, query, projection, sort *bson.M, skip, limit *int64) ([]*T, error) {
	var (
		result []*T
		err    error
	)

	collection := a.Client.GetCollection(a.collectionName)

	options := options.FindOptions{}
	if sort != nil {
		options.Sort = sort
	}
	if projection != nil {
		options.Projection = projection
	}
	if skip != nil {
		options.Skip = skip
	}
	if limit != nil {
		options.Limit = limit
	}

	cursor, err := collection.Find(ctx, query, &options)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		result = []*T{}
	}

	return result, err
}

// Find find slice of entries by filter values
func (r *Repository[T]) FindOne(ctx context.Context, query bson.M) (*T, error) {
	var (
		res T
		err error
	)

	collection := r.Client.GetCollection(r.collectionName)
	result := collection.FindOne(ctx, query)

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMongo, err.Error())
	}

	if err = result.Decode(&res); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResult, err.Error())
	}

	return &res, nil
}

// InsertOne persist entry and returns the primary key
func (r *Repository[T]) InsertOne(ctx context.Context, data *T) (string, error) {
	collection := r.Client.GetCollection(r.collectionName)

	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrMongo, err.Error())
	}

	return result.InsertedID.(string), nil
}

// UpdateOne updates entry that matches id with the value update
func (r *Repository[_]) UpdateOne(ctx context.Context, id string, update interface{}) error {
	if _, err := r.Client.GetCollection(r.collectionName).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update}); err != nil {
		return fmt.Errorf("%w: %s", ErrMongo, err.Error())
	}

	return nil
}

func (r *Repository[_]) Inc(ctx context.Context, fieldName string) error {
	if _, err := r.Client.GetCollection(r.collectionName).UpdateOne(ctx, bson.M{}, bson.M{"$inc": bson.M{fieldName: 1}}); err != nil {
		return fmt.Errorf("%w: %s", ErrMongo, err.Error())
	}

	return nil
}

func (r *Repository[T]) FindByAggregate(ctx context.Context, aggregate interface{}) ([]*T, error) {
	var (
		result []*T
		err    error
	)

	collection := r.Client.GetCollection(r.collectionName)

	// Execute aggregation
	cursor, err := collection.Aggregate(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
