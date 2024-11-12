package apartment

import (
	application "fewoserv/internal/application/apartment"
	applicationPicture "fewoserv/internal/application/picture"
	applicationProcessLog "fewoserv/internal/application/process_log"
	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all real_estate endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processlog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processlog)
	pictureApp := applicationPicture.New(mongoDBClient, processlog)

	handler := NewHandler(app, pictureApp)

	router.HandleFunc("/apartments", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/apartments/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/apartments/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/apartments/{id}/ro", handler.GetReadOnly).Methods(http.MethodGet)
	router.HandleFunc("/apartments/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/apartments", handler.GetMany).Methods(http.MethodGet)
}
