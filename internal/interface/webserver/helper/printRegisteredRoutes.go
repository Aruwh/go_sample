package helper

import (
	"fewoserv/internal/infrastructure/logger"
	"fmt"

	"github.com/gorilla/mux"
)

var log = logger.New("WEBSERVER")

// listRegisteredRoutes is used to list all registered routes
func ListRegisteredRoutes(router *mux.Router) error {
	if err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		registeredRoute, _ := route.GetPathTemplate()
		usedMethod, _ := route.GetMethods()

		if len(usedMethod) != 0 {
			log.Info(fmt.Sprintf("registered route %s %s", usedMethod[0], registeredRoute))
		}

		return nil
	}); err != nil {
		log.Error(fmt.Sprintf("failed to register routes %s", err))
		return err
	}

	return nil
}
