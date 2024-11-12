package attribute

import (
	application "fewoserv/internal/application/attribute"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strconv"
)

var (
	DOMAIN = "Attribute"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Create(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetMany(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
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

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData CreateRequest
	err = helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	transformedName := shared.TransformRequestTranslation(requestData.Name)

	attribute, err := h.application.Create(identity.UserID, transformedName, requestData.Svg)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, attribute)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_DELETE)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	recordID := helper.GetQueryVar(r, "id")

	err = h.application.Delete(identity.UserID, *recordID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, shared.BuildDeleteResponse(true))
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")

	attribute, err := h.application.Get(*recordID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, attribute)
}

func (h Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_VIEW)
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

	attribute, err := h.application.GetMany(requestData.Filter.Name, requestData.Sort, requestData.Skip, requestData.Limit)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, attribute)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData UpdateRequest
	err = helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")

	transformedName := shared.TransformRequestTranslation(requestData.Name)

	attribute, err := h.application.Update(identity.UserID, *recordID, transformedName, requestData.Svg)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, attribute)
}
