package repository

import (
	"context"
	process_log "fewoserv/internal/domain/process_log"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "process_log"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository *mongodb.Repository[process_log.ProcessLog]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[process_log.ProcessLog](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) Insert(process_log *process_log.ProcessLog) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, process_log)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindMany(userID string) ([]*process_log.ProcessLog, error) {
	ctx := context.Background()

	var (
		query       = bson.M{"userID": userID}
		limit int64 = 100
	)

	sort := common.Sort{Order: common.OrderDESC, Field: common.SortByCreated}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), nil, &limit)
}
