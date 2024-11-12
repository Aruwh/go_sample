package permission

import (
	application "fewoserv/internal/application/permission"
	applicationProcessLog "fewoserv/internal/application/process_log"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all permission endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient) {
	processLog := applicationProcessLog.New(mongoDBClient)

	app := application.New(mongoDBClient, processLog)
	handler := NewHandler(app)

	router.HandleFunc("/permissions", handler.GetMany).Methods(http.MethodGet)
}
