package repository

import (
	"context"
	"fewoserv/internal/domain/settings"
	"fewoserv/internal/domain/shared"

	"fewoserv/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "settings"
	cachedRepo     *Repo = nil
)

type (
	Repo struct {
		repository *mongodb.Repository[settings.Settings]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository: mongodb.NewRepository[settings.Settings](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) Insert(settings *settings.Settings) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, settings)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetBookingNumber() (int, error) {
	ctx := context.Background()

	query := bson.M{}
	foundSettings, err := r.repository.Find(ctx, &query, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}
	if len(foundSettings) == 0 {
		return 0, nil
	}

	return foundSettings[0].BookingNumber, nil
}

func (r *Repo) GetNotificationMessage() (*shared.Translation, error) {
	ctx := context.Background()

	query := bson.M{}
	foundSettings, err := r.repository.Find(ctx, &query, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(foundSettings) == 0 {
		return nil, nil
	}

	return foundSettings[0].NotificationMessage, nil
}

func (r *Repo) IncBookingNumber() error {
	ctx := context.Background()

	return r.repository.Inc(ctx, "bookingNumber")
}
