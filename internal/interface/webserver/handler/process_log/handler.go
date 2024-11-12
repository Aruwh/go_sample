package saison

import (
	application "fewoserv/internal/application/process_log"
	"fewoserv/internal/infrastructure/common"
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
		GetMany(w http.ResponseWriter, r *http.Request)
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

func (h Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	adminUserID := helper.GetQueryVar(r, "id")

	message, err := h.application.GetManyBy(*adminUserID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, message)
}
