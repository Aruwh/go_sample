package user

import (
	applicationProcessLog "fewoserv/internal/application/process_log"
	application "fewoserv/internal/application/user"
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all healthz endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, jwtExpireTimeForPwdResetInMinutes int, landingpageEndpoint string) {
	processLog := applicationProcessLog.New(mongoDBClient)
	app := application.New(mongoDBClient, emailHandler, processLog, &landingpageEndpoint, &jwtExpireTimeForPwdResetInMinutes)
	handler := NewHandler(app)

	router.HandleFunc("/users/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/users", handler.GetMany).Methods(http.MethodGet)
	router.HandleFunc("/users", handler.Create).Methods(http.MethodPost)
	// router.HandleFunc("/users/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}", handler.Delete).Methods(http.MethodDelete)
}
