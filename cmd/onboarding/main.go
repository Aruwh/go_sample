package main

import (
	"context"
	applicationAdminUser "fewoserv/internal/application/admin_user"
	applicationAttribute "fewoserv/internal/application/attribute"
	applicationPermisison "fewoserv/internal/application/permission"
	applicationProcessLog "fewoserv/internal/application/process_log"
	"fewoserv/internal/infrastructure/config"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

var log = logger.New("ON_BOARDING")

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

func main() {
	cfg := config.Load()

	mongoDbClient := setupMongoDB(cfg.MongoDB.MongoDBName, cfg.MongoDB.MongoDBUri)
	connectMongoDB(mongoDbClient)

	appProcessLog := applicationProcessLog.New(mongoDbClient)
	appAdminUser := applicationAdminUser.New(mongoDbClient, nil, appProcessLog, nil, nil)
	appPermission := applicationPermisison.New(mongoDbClient, appProcessLog)
	appAttribute := applicationAttribute.New(mongoDbClient, appProcessLog)

	appPermission.CreateDefaults()

	superAdminUser, err := appAdminUser.CreateSuperAdmin()
	if err != nil {
		log.Panic(err.Error())
	}

	appAttribute.CreateDefaults(superAdminUser.ID)

	shutdownMongoDBConnection(mongoDbClient)
}
