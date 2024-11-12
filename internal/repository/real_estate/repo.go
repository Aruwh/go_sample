package repository

import (
	"context"
	realEstate "fewoserv/internal/domain/real_estate"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collectionName       = "real_estate"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[realEstate.RealEstate]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[realEstate.RealEstate](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*realEstate.RealEstate, error) {
	// we use the get many func to ensure to get the linked apartments as well
	records, err := r.FindMany([]*string{&id}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return records[0], nil
}

func (r *Repo) ValidateRealEstateExists(realEstate *realEstate.RealEstate) bool {
	ctx := context.Background()

	query := bson.M{"name": realEstate.Name}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Insert(realEstate *realEstate.RealEstate) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, realEstate)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, realEstate *realEstate.RealEstate) error {
	ctx := context.Background()

	realEstate.Edited.Update(updaterID)
	return r.repository.UpdateOne(ctx, realEstate.ID, realEstate)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(ids []*string, name *string, sort *common.Sort, skip, limit *int64) ([]*realEstate.RealEstate, error) {
	ctx := context.Background()

	query := bson.M{}
	if name != nil {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if ids != nil {
		query["_id"] = bson.M{"$in": ids}
	}

	pipeline := []bson.M{
		{"$match": query},
		{
			"$lookup": bson.M{
				"from":         "apartment",
				"localField":   "_id",
				"foreignField": "realEstateID",
				"as":           "apartments",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$apartments",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$group": bson.M{
				"_id":         "$_id",
				"name":        bson.M{"$first": "$name"},
				"description": bson.M{"$first": "$description"},
				"pictureID":   bson.M{"$first": "$pictureID"},
				"created":     bson.M{"$first": "$created"},
				"edited":      bson.M{"$first": "$edited"},
				"apartments": bson.M{
					"$push": bson.M{
						"_id":  "$apartments._id",
						"name": "$apartments.name",
					},
				},
			},
		},
	}

	var usedSort *primitive.M
	if sort != nil {
		usedSort = sort.ToBson()
		pipeline = append(pipeline, bson.M{"$sort": *usedSort})
	}

	return r.repository.FindByAggregate(ctx, pipeline)
}
