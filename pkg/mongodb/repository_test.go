package mongodb_test

import (
	"context"
	adminuser "fewoserv/internal/domain/admin_user"
	"fewoserv/internal/infrastructure/config"
	"fewoserv/pkg/mongodb"
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func setupMongoDB(mongoDBName, mongoDBUri string) mongodb.IClient {
	dbClient, err := mongodb.NewClient(mongoDBUri, mongoDBName)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to create db client:%s", err))
	}

	return dbClient
}

func connectMongoDB(dbClient mongodb.IClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := dbClient.Connect(ctx)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to connect to mongoDB: %s", err))
	}
}

func shutdownMongoDBConnection(dbClient mongodb.IClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := dbClient.Disconnect(ctx)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to disconnect mongoDB: %s", err))
	}
}

func TestFind(t *testing.T) {
	cfg := config.Load()

	mongoDbClient := setupMongoDB(cfg.MongoDB.MongoDBName, cfg.MongoDB.MongoDBUri)
	connectMongoDB(mongoDbClient)
	repo := mongodb.NewRepository[adminuser.AdminUser](mongoDbClient, "admin_users")

	res, _ := repo.Find(context.Background(), &bson.M{}, nil, nil, nil, nil)

	fmt.Printf("%v,", res)

	shutdownMongoDBConnection(mongoDbClient)
}
