package realestate

import (
	applicationPicture "fewoserv/internal/application/picture"
	applicationProcessLog "fewoserv/internal/application/process_log"
	application "fewoserv/internal/application/real_estate"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all real_estate endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	appPicture := applicationPicture.New(mongoDBClient, processLog)
	handler := NewHandler(app, appPicture)

	router.HandleFunc("/realEstates", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/realEstates/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/realEstates/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/realEstates/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/realEstates", handler.GetMany).Methods(http.MethodGet)
}
