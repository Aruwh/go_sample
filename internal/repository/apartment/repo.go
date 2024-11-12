package repository

import (
	"context"
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collectionName       = "apartment"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository                  *mongodb.Repository[apartment.Apartment]
		repositoryMinimalApartment  *mongodb.Repository[apartment.MinimalApartment]
		repositoryReadOnlyApartment *mongodb.Repository[apartment.ReadOnlyApartment]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository:                  mongodb.NewRepository[apartment.Apartment](dbClient, collectionName),
			repositoryMinimalApartment:  mongodb.NewRepository[apartment.MinimalApartment](dbClient, collectionName),
			repositoryReadOnlyApartment: mongodb.NewRepository[apartment.ReadOnlyApartment](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*apartment.Apartment, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) ValidateRecordExists(apartment *apartment.Apartment) bool {
	ctx := context.Background()

	query := bson.M{"name": apartment.Name}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Insert(apartment *apartment.Apartment) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, apartment)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, apartment *apartment.Apartment) error {
	ctx := context.Background()

	apartment.Edited.Update(updaterID)
	return r.repository.UpdateOne(ctx, apartment.ID, apartment)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindManyMinimal(ownerID, realEstateID, name *string, sort *common.Sort, skip, limit *int64) ([]*apartment.MinimalApartment, error) {
	ctx := context.Background()

	query := bson.M{}
	if name != nil {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if ownerID != nil {
		query["ownerID"] = ownerID
	}
	if realEstateID != nil {
		query["realEstateID"] = realEstateID
	}

	pipeline := []bson.M{
		{
			"$match": query,
		},
		{
			"$lookup": bson.M{
				"from":         "real_estate",
				"localField":   "realEstateID",
				"foreignField": "_id",
				"as":           "realEstate",
			},
		},
		{
			"$unwind": "$realEstate",
		},
		{
			"$addFields": bson.M{
				"realEstateName": "$realEstate.name",
			},
		},
		{
			"$project": bson.M{
				"realEstateID":   1,
				"isActive":       1,
				"name":           1,
				"description":    1,
				"pictureIDs":     1,
				"realEstateName": 1,
				"_id":            1,
			},
		},
	}

	var usedSort *primitive.M
	if sort != nil {
		usedSort = sort.ToBson()
		pipeline = append(pipeline, bson.M{"$sort": *usedSort})
	}

	return r.repositoryMinimalApartment.FindByAggregate(ctx, pipeline)
}

func (r *Repo) FindManyPublic(ownerID, realEstateID, name *string, sort *common.Sort, skip, limit *int64) ([]*apartment.Apartment, error) {
	ctx := context.Background()

	query := bson.M{"isActive": true}
	if name != nil {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if ownerID != nil {
		query["ownerID"] = ownerID
	}
	if realEstateID != nil {
		query["realEstateID"] = realEstateID
	}

	var usedSort *primitive.M
	if sort != nil {
		usedSort = sort.ToBson()
	}

	return r.repository.Find(ctx, &query, nil, usedSort, skip, limit)
}

func (r *Repo) GetReadOnly(recordID, ownerID string) (*apartment.ReadOnlyApartment, error) {
	ctx := context.Background()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": recordID,
				// "ownerID": ownerID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "attribute",
				"localField":   "attributeCollection.topAttributeIDs",
				"foreignField": "_id",
				"as":           "topAttributes",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "attribute",
				"localField":   "attributeCollection.attributeIDs",
				"foreignField": "_id",
				"as":           "attributes",
			},
		},
		{
			"$project": bson.M{
				"name":        1,
				"description": 1,
				"pictureIDs":  1,
				"attributeCollection": bson.M{
					"roomSize":              1,
					"sleepingPlaces":        1,
					"allowedNumberOfPeople": 1,
					"topAttributes": bson.M{
						"$map": bson.M{
							"input": "$topAttributes",
							"as":    "topAttr",
							"in": bson.M{
								"name": "$$topAttr.name",
								"svg":  "$$topAttr.svg",
							},
						},
					},
					"attributes": bson.M{
						"$map": bson.M{
							"input": "$attributes",
							"as":    "attr",
							"in": bson.M{
								"name": "$$attr.name",
								"svg":  "$$attr.svg",
							},
						},
					},
				},
				"saisonPrice": bson.M{
					"lowPrice":    1,
					"middlePrice": 1,
					"highPrice":   1,
					"peakPrice":   1,
				},
			},
		},
	}

	records, err := r.repositoryReadOnlyApartment.FindByAggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return records[0], nil
}
