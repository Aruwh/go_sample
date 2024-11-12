package healthz

import (
	"github.com/gorilla/mux"
)

// RegisterRouter registered all healthz endpoints
func RegisterRouter(router *mux.Router) {
	handler := NewHandler()

	router.HandleFunc("/healthz", handler.Healthz).Methods("GET")
	router.HandleFunc("/ping", handler.Ping).Methods("GET")
}
