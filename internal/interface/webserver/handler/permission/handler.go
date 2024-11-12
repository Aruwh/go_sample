package permission

import (
	application "fewoserv/internal/application/permission"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strconv"
)

var (
	DOMAIN = "Permission"
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
func NewHandler(applicationAdminUser application.IApplication) IHandler {
	var handler = Handler{}

	handler.response.Domain = DOMAIN
	handler.application = applicationAdminUser

	return handler
}

func (h Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PERMISSION_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData common.GetManyRequest[GetManyFilter]
	name := helper.GetQueryParam(r, "name")
	requestData.Filter.Name = name

	field := helper.GetQueryParam(r, "field")
	requestData.Sort.Field = common.SortByType(*field)

	order := helper.GetQueryParam(r, "order")
	requestData.Sort.Order = common.OrderType(*order)

	limit, _ := strconv.ParseInt(*helper.GetQueryParam(r, "limit"), 10, 64)
	requestData.Limit = limit

	skip, _ := strconv.ParseInt(*helper.GetQueryParam(r, "skip"), 10, 64)
	requestData.Skip = skip

	permission, err := h.application.GetMany(requestData.Filter.Name, requestData.Sort, requestData.Skip, requestData.Limit)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, permission)
}
