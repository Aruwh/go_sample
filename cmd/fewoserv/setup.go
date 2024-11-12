package main

import (
	picturecache "fewoserv/internal/infrastructure/picture_cache"
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/internal/interface/webserver"
	"fewoserv/pkg/mongodb"
	"fmt"
)

func setupWebServer(httpPort, storagePath, feEndpoint, landingpageEndpoint, sslCertPath, sslKeyPath string, jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes int, mongoDBClient mongodb.IClient, corsAllowedOrigins, copyEmailAddresses []string, emialHandler emailhandler.IEmailHandler, pictureCache picturecache.IPictureCache, corsDebuggindEnabled bool) webserver.IWebserver {
	ws, err := webserver.New(httpPort, storagePath, feEndpoint, landingpageEndpoint, sslCertPath, sslKeyPath, jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes, mongoDBClient, emialHandler, pictureCache, corsAllowedOrigins, copyEmailAddresses, corsDebuggindEnabled)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to start web server: %s", err))
	}

	return ws
}

func setupMongoDB(mongoDBName, mongoDBUri string) mongodb.IClient {
	dbClient, err := mongodb.NewClient(mongoDBUri, mongoDBName)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to create db client:%s", err))
	}

	return dbClient
}

func setupEmailHandler(fromEmail, password, smtpServerUrl, smtpServerPort string) emailhandler.IEmailHandler {
	emailHandler := emailhandler.New(fromEmail, password, smtpServerUrl, smtpServerPort)
	return emailHandler
}

func setupPictureCache() picturecache.IPictureCache {
	return picturecache.New()
}
