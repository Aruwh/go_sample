package news

import (
	application "fewoserv/internal/application/news"
	applicationProcessLog "fewoserv/internal/application/process_log"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all real_estate endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	handler := NewHandler(app)

	router.HandleFunc("/news", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/news/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/news/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/news/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/news", handler.GetMany).Methods(http.MethodGet)
}
