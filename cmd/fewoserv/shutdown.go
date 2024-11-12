package main

import (
	"context"
	"fewoserv/internal/interface/webserver"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

func shutdownWebserver(webserver webserver.IWebserver) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := webserver.Stop(ctx)
	if err != nil {
		log.Panic(fmt.Sprintf("unable to shut down web server: %s", err))
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
