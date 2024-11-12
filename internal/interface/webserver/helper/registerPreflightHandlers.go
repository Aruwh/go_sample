package helper

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent)
}

// listRegisteredRoutes is used to list all registered routes
func RegisterPreflightHandlers(router *mux.Router) error {
	if err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		for _, method := range methods {
			if method == "OPTIONS" {
				continue
			}
			router.HandleFunc(pathTemplate, optionsHandler).Methods("OPTIONS")
			log.Info(fmt.Sprintf("preflight handler %s registered", pathTemplate))

		}

		return nil
	}); err != nil {
		log.Error(fmt.Sprintf("failed to register routes %s", err))
		return err
	}

	return nil
}
