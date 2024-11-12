package main

import (
	"context"
	"errors"
	"fewoserv/internal/interface/webserver"
	"fewoserv/pkg/mongodb"
	"fmt"
	"net/http"
	"time"
)

func ServerStartedPrompt() {
	mow := `	         
		                      ,     ,
		                  ___('-&&&-')__
		                 '.__./     \__.'
		     _     _     _ .-'  0  0 \
		   /' '--'( ('--' '\         |
		  /        ) )      \ \ _   _|
		 |        ( (        | (0_._0)
		 |         ) )       |/ '---'
		 |        ( (        |\_
		 |         ) )       |( \,
		  \       ((       / )__/
		   |     /:))\     |   d
		   |    /:((::\    |
		   |   |:::):::|   |
		   /   \::&&:::/   \
		   \   /;U&::U;\   /
		    | | | u:u | | |
		    | | \     / | |
		    | | _|   | _| |
		    / \""'   '""/ \
		   | __|       | __|
		   '"""'       '"""'

                 FeWoServ started !




`

	fmt.Print(mow)
}

func startWebServer(ws webserver.IWebserver) {
	if err := ws.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic(fmt.Sprintf("failed to start http server: %s", err))
	}
}

func connectMongoDB(dbClient mongodb.IClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := dbClient.Connect(ctx)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to connect to mongoDB: %s", err))
	}
}
