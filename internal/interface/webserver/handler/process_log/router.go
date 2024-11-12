package saison

import (
	application "fewoserv/internal/application/process_log"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all picture endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	app := application.New(mongoDBClient)
	handler := NewHandler(app)

	router.HandleFunc("/processLog/{id}", handler.GetMany).Methods(http.MethodGet)
}
