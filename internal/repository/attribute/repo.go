package repository

import (
	"context"
	attribute "fewoserv/internal/domain/attribute"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "attribute"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[attribute.Attribute]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[attribute.Attribute](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*attribute.Attribute, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) ValidateRecordExists(attribute *attribute.Attribute) bool {
	ctx := context.Background()

	query := bson.M{"name": attribute.Name}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Insert(attribute *attribute.Attribute) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, attribute)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, attribute *attribute.Attribute) error {
	ctx := context.Background()

	attribute.Edited.Update(updaterID)
	return r.repository.UpdateOne(ctx, attribute.ID, attribute)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(name *string, sort *common.Sort, skip, limit int64) ([]*attribute.Attribute, error) {
	ctx := context.Background()

	var query = bson.M{}
	if name != nil {
		query["$or"] = []bson.M{
			{"name.deDE": bson.M{"$regex": name, "$options": "i"}},
			{"name.enGB": bson.M{"$regex": name, "$options": "i"}},
			{"name.frFR": bson.M{"$regex": name, "$options": "i"}},
			{"name.itIT": bson.M{"$regex": name, "$options": "i"}},
		}
	}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), &skip, &limit)
}

func (r *Repo) FindManyByID(ids []string) ([]*attribute.Attribute, error) {
	ctx := context.Background()

	var query = bson.M{"_id": bson.M{"$in": ids}}

	return r.repository.Find(ctx, &query, nil, nil, nil, nil)
}
