package repository

import (
	"context"
	adminuser "fewoserv/internal/domain/admin_user"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "admin_users"
	cachedRepo     *Repo = nil
)

type (
	Repo struct {
		repository *mongodb.Repository[adminuser.AdminUser]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[adminuser.AdminUser](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) ValidateAdminUserExists(adminUser *adminuser.AdminUser) bool {
	foundRecord, err := r.LoadAdminUserByEmail(*adminUser.Email)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) LoadAdminUserByEmail(email string) (*adminuser.AdminUser, error) {
	ctx := context.Background()

	query := bson.M{"email": email}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) LoadAdminUserByID(id string) (*adminuser.AdminUser, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) DeleteAdminUserByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(adminUserType *common.AdminUserType, name *string, sort *common.Sort, skip, limit int64) ([]*adminuser.AdminUser, error) {
	ctx := context.Background()

	query := bson.M{"type": bson.M{"$ne": common.SUPER_ADMINISTRATOR}}

	if adminUserType != nil {
		query["type"] = adminUserType
	}

	if name != nil {
		query["$or"] = []bson.M{
			{"firstName": bson.M{"$regex": name, "$options": "i"}},
			{"lastName": bson.M{"$regex": name, "$options": "i"}},
		}
	}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), &skip, &limit)
}

func (r *Repo) Insert(adminUser *adminuser.AdminUser) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, adminUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, adminUser *adminuser.AdminUser) error {
	ctx := context.Background()

	adminUser.Edited.Update(updaterID)

	return r.repository.UpdateOne(ctx, adminUser.ID, adminUser)
}
