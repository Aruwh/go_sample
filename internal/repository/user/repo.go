package repository

import (
	"context"
	user "fewoserv/internal/domain/user"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "users"
	cachedRepo     *Repo = nil
)

type (
	Repo struct {
		repository *mongodb.Repository[user.User]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[user.User](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadUserByEmail(email string) (*user.User, error) {
	ctx := context.Background()

	query := bson.M{"email": email}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) LoadUserByID(id string) (*user.User, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) DeleteUserByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(name *string, sort *common.Sort, skip, limit int64) ([]*user.User, error) {
	ctx := context.Background()

	query := bson.M{}
	if name != nil {
		query["$or"] = []bson.M{
			{"firstName": bson.M{"$regex": name, "$options": "i"}},
			{"lastName": bson.M{"$regex": name, "$options": "i"}},
		}
	}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), &skip, &limit)
}

func (r *Repo) Insert(adminUser *user.User) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, adminUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, adminUser *user.User) error {
	ctx := context.Background()

	adminUser.Edited.Update(updaterID)

	return r.repository.UpdateOne(ctx, adminUser.ID, adminUser)
}
