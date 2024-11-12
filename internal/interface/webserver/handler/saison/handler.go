package saison

import (
	application "fewoserv/internal/application/saison"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strconv"
)

var (
	DOMAIN = "Season"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Create(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetMany(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
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

	transformedSaisonEntries := shared.TransformRequestSaisonEntries(&requestData.Enries)

	attribute, err := h.application.Create(identity.UserID, requestData.Year, *transformedSaisonEntries)
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

	saison, err := h.application.Get(*recordID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, saison)
}

func (h Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.SETTINGS_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData common.GetManyRequest[GetManyFilter]

	searchYear := helper.GetQueryParam(r, "year")
	if searchYear != nil {
		transformed, err := strconv.Atoi(*searchYear)
		if err == nil {
			requestData.Filter.Year = &transformed
		}

	}

	field := helper.GetQueryParam(r, "field")
	requestData.Sort.Field = common.SortByType(*field)

	order := helper.GetQueryParam(r, "order")
	requestData.Sort.Order = common.OrderType(*order)

	limit, _ := strconv.ParseInt(*helper.GetQueryParam(r, "limit"), 10, 64)
	requestData.Limit = limit

	skip, _ := strconv.ParseInt(*helper.GetQueryParam(r, "skip"), 10, 64)
	requestData.Skip = skip

	saisons, err := h.application.GetMany(requestData.Sort, requestData.Skip, requestData.Limit, requestData.Filter.Year)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, saisons)
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

	transformedSaisonEntries := shared.TransformRequestSaisonEntries(requestData.Enries)

	saison, err := h.application.Update(identity.UserID, *recordID, requestData.Year, transformedSaisonEntries)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, saison)
}
