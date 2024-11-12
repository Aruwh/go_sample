package adminuser

import (
	application "fewoserv/internal/application/admin_user"
	applicationProcessLog "fewoserv/internal/application/process_log"
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all healthz endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, jwtExpireTimeForPwdResetInMinutes int, feEndpoint string) {
	processLog := applicationProcessLog.New(mongoDBClient)
	app := application.New(mongoDBClient, emailHandler, processLog, &feEndpoint, &jwtExpireTimeForPwdResetInMinutes)
	handler := NewHandler(app)

	router.HandleFunc("/adminUsers/me", handler.GetMe).Methods(http.MethodGet)
	router.HandleFunc("/adminUsers/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/adminUsers", handler.GetMany).Methods(http.MethodGet)
	router.HandleFunc("/adminUsers", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/adminUsers/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/adminUsers/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/adminUsers/{id}/invite", handler.Invite).Methods(http.MethodPatch)
}
