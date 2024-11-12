package saison

import (
	applicationProcessLog "fewoserv/internal/application/process_log"
	application "fewoserv/internal/application/settings"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all picture endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	handler := NewHandler(app)

	router.HandleFunc("/settings/notification", handler.Get).Methods(http.MethodGet)
}
