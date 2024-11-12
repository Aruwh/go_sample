package repository

import (
	"context"
	"fewoserv/internal/domain/permission"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "permissions"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[permission.Permission]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[permission.Permission](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) ValidatePermissionExists(permission *permission.Permission) bool {
	ctx := context.Background()

	query := bson.M{"name": permission.Name}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Insert(permission *permission.Permission) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, permission)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindMany(name *string, sort *common.Sort, skip, limit int64) ([]*permission.Permission, error) {
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
