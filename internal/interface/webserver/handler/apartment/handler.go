package apartment

import (
	application "fewoserv/internal/application/apartment"
	applicationPicture "fewoserv/internal/application/picture"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strconv"
)

var (
	DOMAIN = "Apartment"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Create(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetReadOnly(w http.ResponseWriter, r *http.Request)
		GetMany(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)

		GetManyPublic(w http.ResponseWriter, r *http.Request)
	}

	// Handler defines the used handler which combines all implementations of available endpoints
	Handler struct {
		response           shared.Response
		application        application.IApplication
		applicationPicture applicationPicture.IApplication
	}
)

// NewHandler creates a new Handler
func NewHandler(applicationAdminUser application.IApplication, applicationPicture applicationPicture.IApplication) IHandler {
	var handler = Handler{}

	handler.response.Domain = DOMAIN
	handler.application = applicationAdminUser
	handler.applicationPicture = applicationPicture

	return handler
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.APARTMENT_EDIT)
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

	transformedDescription := shared.TransformRequestTranslation(requestData.Description)

	apartment, err := h.application.Create(identity.UserID, requestData.OwnerID, requestData.RealEstateID, requestData.Name, transformedDescription, requestData.PictureIDs, requestData.AttributeIDs, requestData.TopAttributeIDs, requestData.SaisonPrice, requestData.RoomSize, requestData.SleepingPlaces, requestData.Bathrooms, requestData.AllowedNumberOfPeople, requestData.AllowedNumberOfPets)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartment)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.APARTMENT_DELETE)
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
	var (
		apartment interface{}
		err       error
	)

	identity := helper.ExtractIdentity(r)
	err = identity.EnsureRequestPermission(common.APARTMENT_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")
	apartment, err = h.application.Get(*recordID)

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartment)
}

func (h Handler) GetReadOnly(w http.ResponseWriter, r *http.Request) {
	var (
		apartment interface{}
		err       error
	)

	identity := helper.ExtractIdentity(r)
	err = identity.EnsureRequestPermission(common.APARTMENT_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")
	apartment, err = h.application.GetReadOnly(*recordID, identity.UserID)

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartment)
}

func (h Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.APARTMENT_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData common.GetManyRequest[GetManyFilter]

	name := helper.GetQueryParam(r, "name")
	requestData.Filter.Name = name

	realEstateID := helper.GetQueryParam(r, "realEstateID")
	requestData.Filter.RealEstateID = realEstateID

	field := helper.GetQueryParam(r, "field")
	requestData.Sort.Field = common.SortByType(*field)

	order := helper.GetQueryParam(r, "order")
	requestData.Sort.Order = common.OrderType(*order)

	limit, _ := strconv.ParseInt(*helper.GetQueryParam(r, "limit"), 10, 64)
	requestData.Limit = limit

	skip, _ := strconv.ParseInt(*helper.GetQueryParam(r, "skip"), 10, 64)
	requestData.Skip = skip

	var userID *string
	if identity.Type == common.AUDIANCE_APARTMENT_OWNER {
		userID = &identity.UserID
	}

	apartments, err := h.application.GetMany(userID, requestData.Filter.RealEstateID, requestData.Filter.Name, &requestData.Sort, &requestData.Skip, &requestData.Limit)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartments)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.APARTMENT_EDIT)
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

	transformedDescription := shared.TransformRequestTranslation(requestData.Description)

	apartment, pictureIDsToRemove, err := h.application.Update(identity.UserID, *recordID, requestData.IsActive, requestData.OwnerID, requestData.RealEstateID, requestData.Name, transformedDescription, requestData.PictureIDs, requestData.AttributeIDs, requestData.TopAttributeIDs, requestData.SaisonPrice, requestData.RoomSize, requestData.Bathrooms, requestData.SleepingPlaces, requestData.AllowedNumberOfPeople, requestData.AllowedNumberOfPets)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	for _, pictureIDToRemove := range pictureIDsToRemove {
		h.applicationPicture.Delete(identity.UserID, pictureIDToRemove)
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartment)
}

func (h Handler) GetManyPublic(w http.ResponseWriter, r *http.Request) {
	apartments, err := h.application.GetManyPublic()
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, NewGetManyPublicResponse(apartments))
}
