package saison

import (
	application "fewoserv/internal/application/settings"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
)

var (
	DOMAIN = "Settings"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Get(w http.ResponseWriter, r *http.Request)
	}

	// Handler defines the used handler which combines all implementations of available endpoints
	Handler struct {
		response    shared.Response
		application application.IApplication
	}
)

// NewHandler creates a new Handler
func NewHandler(applicationPicture application.IApplication) IHandler {
	var handler = Handler{}

	handler.response.Domain = DOMAIN
	handler.application = applicationPicture

	return handler
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission()
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	message, err := h.application.GetNotificationMessage()
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, message)
}
