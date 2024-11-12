package main

import (
	"fewoserv/internal/infrastructure/config"
	"fewoserv/internal/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
)

var log = logger.New("MAIN")

func signalChannel() chan os.Signal {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	return done
}

func main() {
	cfg := config.Load()

	emailHandler := setupEmailHandler(cfg.Email.From, cfg.Email.Password, cfg.Email.ServerURL, cfg.Email.ServerPort)

	mongoDbClient := setupMongoDB(cfg.MongoDB.MongoDBName, cfg.MongoDB.MongoDBUri)
	connectMongoDB(mongoDbClient)

	pictureCache := setupPictureCache()

	webserver := setupWebServer(
		cfg.Service.HTTPPort,
		cfg.Service.StoragePath,
		cfg.Email.FeEndpoint,
		cfg.Email.LandingpageEndpoint,
		cfg.Ssl.CertPath,
		cfg.Ssl.KeyPath,
		cfg.Service.JwtExpireTimeInMinutes,
		cfg.Authentication.JwtExpireTimeForPwdResetInMinutes,
		mongoDbClient,
		cfg.Service.CORSAllowedOrigins,
		cfg.Email.CopyEmailAddresses,
		emailHandler,
		pictureCache,
		cfg.Service.CORS_DEBUGING)

	go startWebServer(webserver)

	ServerStartedPrompt()

	<-signalChannel()

	shutdownMongoDBConnection(mongoDbClient)
	shutdownWebserver(webserver)
}
