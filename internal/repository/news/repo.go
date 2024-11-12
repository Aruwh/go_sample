package repository

import (
	"context"
	news "fewoserv/internal/domain/news"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "news"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[news.News]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[news.News](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*news.News, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) Insert(news *news.News) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, news)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, news *news.News) error {
	ctx := context.Background()

	news.Edited.Update(updaterID)
	return r.repository.UpdateOne(ctx, news.ID, news)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(title *string, sort *common.Sort, skip, limit int64) ([]*news.News, error) {
	ctx := context.Background()

	var query = bson.M{}
	if title != nil {
		query["$or"] = []bson.M{
			{"title.deDE": bson.M{"$regex": title, "$options": "i"}},
			{"title.enGB": bson.M{"$regex": title, "$options": "i"}},
			{"title.frFR": bson.M{"$regex": title, "$options": "i"}},
			{"title.itIT": bson.M{"$regex": title, "$options": "i"}},
		}
	}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), &skip, &limit)
}

func (r *Repo) FindManyPublic() ([]*news.News, error) {
	ctx := context.Background()

	var query = bson.M{
		"active":    true,
		"publishAt": bson.M{"$lte": time.Now()},
	}

	sort := common.Sort{Order: common.OrderDESC, Field: common.SortByCreated}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), nil, nil)
}
