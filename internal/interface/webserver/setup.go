package webserver

import (
	"crypto/tls"
	handlerAdminUser "fewoserv/internal/interface/webserver/handler/admin_user"
	handlerApartment "fewoserv/internal/interface/webserver/handler/apartment"
	handlerAttribute "fewoserv/internal/interface/webserver/handler/attribute"
	handlerAuthentication "fewoserv/internal/interface/webserver/handler/authentication"
	handlerBooking "fewoserv/internal/interface/webserver/handler/booking"
	handlerHealthz "fewoserv/internal/interface/webserver/handler/healthz"
	handlerNews "fewoserv/internal/interface/webserver/handler/news"
	handlerPermisison "fewoserv/internal/interface/webserver/handler/permission"
	handlerPicture "fewoserv/internal/interface/webserver/handler/picture"
	handlerProcessLogs "fewoserv/internal/interface/webserver/handler/process_log"
	handlerPublic "fewoserv/internal/interface/webserver/handler/public"
	handlerRealEstate "fewoserv/internal/interface/webserver/handler/real_estate"
	handlerSaison "fewoserv/internal/interface/webserver/handler/saison"
	handlerSettings "fewoserv/internal/interface/webserver/handler/settings"
	handlerUser "fewoserv/internal/interface/webserver/handler/user"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// applyMiddlewares adds different kind of middlewears to the internal used server
func (ws *Webserver) applyMiddlewares(router *mux.Router) {
	log.Info("applying middlewares...")

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.CorrelatorMiddleware)
	router.Use(middleware.RequestMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.ResponseMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthenticationMiddleware)
}

func (ws *Webserver) registerHandlerToRouter() *mux.Router {
	router := mux.NewRouter()

	// set health endpoint on root path
	handlerAuthentication.RegisterRouter(router, ws.jwtExpireTimeInMinutes, ws.jwtExpireTimeForPwdResetInMinutes, ws.feEndpoint, ws.mongoDBClient, ws.emailhandler)
	handlerHealthz.RegisterRouter(router)
	handlerAdminUser.RegisterRouter(router, ws.mongoDBClient, ws.emailhandler, ws.jwtExpireTimeForPwdResetInMinutes, ws.feEndpoint)
	handlerRealEstate.RegisterRouter(router, ws.mongoDBClient)
	handlerApartment.RegisterRouter(router, ws.mongoDBClient)
	handlerAttribute.RegisterRouter(router, ws.mongoDBClient)
	handlerBooking.RegisterRouter(router, ws.mongoDBClient, ws.emailhandler, ws.landingpageEndpoint, ws.copyEmailAddresses)
	handlerPicture.RegisterRouter(router, ws.mongoDBClient, ws.storagePath, ws.pictureCache)
	handlerSaison.RegisterRouter(router, ws.mongoDBClient)
	handlerPermisison.RegisterRouter(router, ws.mongoDBClient)
	handlerNews.RegisterRouter(router, ws.mongoDBClient)
	handlerSettings.RegisterRouter(router, ws.mongoDBClient)
	handlerProcessLogs.RegisterRouter(router, ws.mongoDBClient)
	handlerPublic.RegisterRouter(router, ws.mongoDBClient, ws.emailhandler, ws.landingpageEndpoint, ws.copyEmailAddresses)
	handlerUser.RegisterRouter(router, ws.mongoDBClient, ws.emailhandler, ws.jwtExpireTimeForPwdResetInMinutes, ws.landingpageEndpoint)

	return router
}

func (ws *Webserver) setupCORSHandling(router *mux.Router) http.Handler {
	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"http://localhost", "http://localhost:3000"},
		AllowedOrigins:   ws.corsAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		//AllowedOrigins:   []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Bearer", "Bearer ", "Origin", "Accept", "Authorization"},

		// Enable Debugging for testing, consider disabling in production
		Debug:                ws.corsDebuggindEnabled,
		OptionsPassthrough:   true,
		OptionsSuccessStatus: 1,
	})

	return c.Handler(router)
}

func (ws *Webserver) generateTLS() (error, *tls.Config) {
	certificate, err := tls.LoadX509KeyPair(ws.sslCertPath, ws.sslKeyPath)
	if err != nil {
		return err, nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	return nil, tlsConfig
}

// configure configures all needed clients and registers all routes
func (ws *Webserver) configure() (error, *mux.Router) {
	log.Info("configuring web server...")

	ws.server.Addr = fmt.Sprintf(":%s", ws.port)
	if ws.sslCertPath != "" && ws.sslKeyPath != "" {
		err, tlsConf := ws.generateTLS()
		if err != nil {
			return err, nil
		}
		ws.server.TLSConfig = tlsConf
	}

	router := ws.registerHandlerToRouter()
	ws.applyMiddlewares(router)

	ws.server.Handler = ws.setupCORSHandling(router)

	err := helper.ListRegisteredRoutes(router)
	if err != nil {
		return err, nil
	}

	err = helper.RegisterPreflightHandlers(router)
	if err != nil {
		return err, nil
	}

	return nil, router
}
