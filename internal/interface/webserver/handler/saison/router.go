package saison

import (
	applicationProcessLog "fewoserv/internal/application/process_log"
	application "fewoserv/internal/application/saison"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all picture endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	handler := NewHandler(app)

	router.HandleFunc("/saisons", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/saisons/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/saisons/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/saisons", handler.GetMany).Methods(http.MethodGet)
	router.HandleFunc("/saisons/{id}", handler.Update).Methods(http.MethodPatch)
}
