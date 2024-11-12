package attribute

import (
	application "fewoserv/internal/application/attribute"
	applicationProcessLog "fewoserv/internal/application/process_log"
	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all real_estate endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processlog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processlog)
	handler := NewHandler(app)

	router.HandleFunc("/attributes", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/attributes/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/attributes/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/attributes/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/attributes", handler.GetMany).Methods(http.MethodGet)
}
