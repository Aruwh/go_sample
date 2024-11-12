package picture

import (
	application "fewoserv/internal/application/picture"
	applicationProcessLog "fewoserv/internal/application/process_log"
	picturecache "fewoserv/internal/infrastructure/picture_cache"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all picture endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient, storagePath string, pictureCache picturecache.IPictureCache) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	handler := NewHandler(app, storagePath, pictureCache)

	router.HandleFunc("/pictures/upload", handler.Upload).Methods(http.MethodPost)
	router.HandleFunc("/pictures/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/pictures/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/pictures/{id}", handler.Update).Methods(http.MethodPatch)
}
