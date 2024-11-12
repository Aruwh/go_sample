package repository

import (
	"context"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collectionName       = "saison"
	cachedRepo     *Repo = nil
)

type (
	Repo struct {
		repository *mongodb.Repository[shared.Saison]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[shared.Saison](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) ValidateRealEstateExists(saison *shared.Saison) (*shared.Saison, error) {
	ctx := context.Background()

	query := bson.M{"year": saison.Year}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) LoadByID(id string) (*shared.Saison, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(sort *common.Sort, skip, limit int64, searchYear *int) ([]*shared.Saison, error) {
	ctx := context.Background()

	query := bson.M{}
	if searchYear != nil {
		query["year"] = searchYear
	}

	var usedSort *primitive.M
	if sort != nil {
		usedSort = sort.ToBson()
	}

	return r.repository.Find(ctx, &query, nil, usedSort, &skip, &limit)
}

func (r *Repo) FindManyByYears(years []int) ([]*shared.Saison, error) {
	ctx := context.Background()

	query := bson.M{}
	query["year"] = bson.M{"$in": years}

	return r.repository.Find(ctx, &query, nil, nil, nil, nil)
}

func (r *Repo) Insert(saison *shared.Saison) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, saison)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, saison *shared.Saison) error {
	ctx := context.Background()

	return r.repository.UpdateOne(ctx, saison.ID, saison)
}
