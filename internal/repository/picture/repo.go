package repository

import (
	"context"
	picture "fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "picture"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[picture.Picture]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[picture.Picture](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*picture.Picture, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) ValidateRecordExists(picture *picture.Picture) bool {
	ctx := context.Background()

	query := bson.M{"_id": picture.ID}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Upsert(picture *picture.Picture) error {
	record, _ := r.LoadByID(picture.ID)
	if record != nil {
		return r.Update(picture)
	}

	return r.Insert(picture)
}

func (r *Repo) Insert(picture *picture.Picture) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, picture)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(picture *picture.Picture) error {
	ctx := context.Background()

	return r.repository.UpdateOne(ctx, picture.ID, picture)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(pictureIDs []string) ([]*picture.Picture, error) {
	ctx := context.Background()

	query := bson.M{}
	query["_id"] = bson.M{"$in": pictureIDs}

	return r.repository.Find(ctx, &query, nil, nil, nil, nil)
}
