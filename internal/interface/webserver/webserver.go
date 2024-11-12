package webserver

import (
	"context"
	"fewoserv/internal/infrastructure/logger"
	picturecache "fewoserv/internal/infrastructure/picture_cache"
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/pkg/mongodb"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var log = logger.New("WEBSERVER")

// Webserver encapsulates the a default go http server.
type (
	IWebserver interface {
		Start() (err error)
		Stop(ctx context.Context) error
	}

	Webserver struct {
		port                              string
		storagePath                       string
		feEndpoint                        string
		landingpageEndpoint               string
		server                            *http.Server
		mongoDBClient                     mongodb.IClient
		emailhandler                      emailhandler.IEmailHandler
		pictureCache                      picturecache.IPictureCache
		startTime                         time.Time
		jwtExpireTimeInMinutes            int
		jwtExpireTimeForPwdResetInMinutes int
		corsAllowedOrigins                []string
		copyEmailAddresses                []string
		corsDebuggindEnabled              bool
		sslCertPath                       string
		sslKeyPath                        string

		router *mux.Router
	}
)

// New creates a new Webserver object with given handlers and returns an pointer
// on the new webserver.
func New(httpPort, storagePath, feEndpoint, landingpageEndpoint, sslCertPath, sslKeyPath string, jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes int, mongoDbClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, pictureCache picturecache.IPictureCache, corsAllowedOrigins, copyEmailAddresses []string, corsDebuggindEnabled bool) (IWebserver, error) {
	if len(httpPort) == 0 {
		return nil, fmt.Errorf("port must not be empty")
	}

	if len(storagePath) == 0 {
		return nil, fmt.Errorf("storagePath must not be empty")
	}

	if jwtExpireTimeInMinutes <= 10 {
		return nil, fmt.Errorf("jwt expire time needs to be > 10")
	}

	webserver := Webserver{
		port:                httpPort,
		storagePath:         storagePath,
		feEndpoint:          feEndpoint,
		landingpageEndpoint: landingpageEndpoint,
		server: &http.Server{
			// NOTE: Timeouts on the webserver are set to prevent slowloris attacks
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
		},
		startTime:                         time.Now().UTC(),
		mongoDBClient:                     mongoDbClient,
		emailhandler:                      emailHandler,
		pictureCache:                      pictureCache,
		jwtExpireTimeInMinutes:            jwtExpireTimeInMinutes,
		jwtExpireTimeForPwdResetInMinutes: jwtExpireTimeForPwdResetInMinutes,
		corsAllowedOrigins:                corsAllowedOrigins,
		copyEmailAddresses:                copyEmailAddresses,
		corsDebuggindEnabled:              corsDebuggindEnabled,
		sslCertPath:                       sslCertPath,
		sslKeyPath:                        sslKeyPath,
	}

	err, router := webserver.configure()
	if err != nil {
		return nil, err
	}

	webserver.router = router
	return &webserver, nil
}

// Start starts the webserver and opens the port.
func (ws *Webserver) Start() (err error) {
	if ws.sslCertPath != "" && ws.sslKeyPath != "" {
		log.Info(fmt.Sprintf("starting secure web server on port: %s", ws.port))
		
		return http.ListenAndServeTLS(fmt.Sprintf(":%s", ws.port), ws.sslCertPath, ws.sslKeyPath, ws.router)
	}

	log.Info(fmt.Sprintf("starting web server on port: %s", ws.port))

	return ws.server.ListenAndServe()
}

// stop stops the webserver gracefully.
func (ws *Webserver) Stop(ctx context.Context) error {
	log.Info("shutting down web server...")

	return ws.server.Shutdown(ctx)
}
