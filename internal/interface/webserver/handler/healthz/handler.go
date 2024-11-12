package healthz

import (
	"fewoserv/internal/interface/webserver/shared"
	"fmt"
	"net/http"
)

var (
	DOMAIN = "Healthz"
)

type (
	// Handler Interface defines the structure of functions which needed to cover all route implementations
	Handler interface {
		Ping(w http.ResponseWriter, r *http.Request)
		Healthz(w http.ResponseWriter, r *http.Request)
	}

	// HealthHandler defines the used handler which combines all implementations of available endpoints
	HealthHandler struct {
		response shared.Response
	}
)

// NewHandler creates a new HealthHandler
func NewHandler() Handler {
	var handler = HealthHandler{}
	handler.response.Domain = DOMAIN

	return handler
}

// Get implements logic to get health information of the mono api and the searcher service
func (h HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	customData := r.Context().Value("token")
	fmt.Println(customData)

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, Healtzh{}.ConvertFromValue(true))
}

// BuildGetPayload extracts relevant informations from the ctx body/query/params and constructs a proper payload, which will
// be passed to the used client
func (h HealthHandler) BuildGetPayload(_ *http.Request) (map[string]interface{}, error) {
	payload := map[string]interface{}{}

	return payload, nil
}

// Ping implements logic to build a pong request
func (h HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	h.response.CreateSuccess(r.Context(), w, http.StatusOK, Ping{}.ConvertFromValue(true))
}
