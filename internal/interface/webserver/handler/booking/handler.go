package booking

import (
	application "fewoserv/internal/application/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"fmt"
	"net/http"
	"time"
)

var (
	DOMAIN = "Booking"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Create(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetBookingOverviews(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		AddMessage(w http.ResponseWriter, r *http.Request)
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
	err := identity.EnsureRequestPermission(common.BOOKING_EDIT)
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

	booking, err := h.application.Create(identity.UserID, requestData.ApartmentID, nil, requestData.Status, requestData.FromDate, requestData.ToDate, requestData.AdultAmount, requestData.ChildAmount, requestData.PetAmount)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, booking)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.BOOKING_DELETE)
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
	err := identity.EnsureRequestPermission(common.BOOKING_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")

	booking, err := h.application.Get(*recordID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, booking)
}

func (h Handler) GetBookingOverviews(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.BOOKING_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData common.GetManyRequest[GetManyBookingOverviewsRequest]

	dateString := helper.GetQueryParam(r, "date")
	if dateString != nil {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, *dateString)
		if err != nil {
			fmt.Println("Fehler beim Parsen:", err)
			return
		}
		requestData.Filter.Date = parsedDate
	}

	var (
		userID       *string
		audianceType *common.AudianceType
	)

	if identity.Type == common.AUDIANCE_APARTMENT_OWNER {
		userID = &identity.UserID
		audianceType = &identity.Type
	}

	booking, err := h.application.GetBookingOverviews(userID, audianceType, requestData.Filter.Date)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, booking)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.BOOKING_EDIT)
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

	booking, err := h.application.Update(identity.UserID, *recordID, requestData.Status, requestData.FromDate, requestData.ToDate, requestData.AdultAmount, requestData.ChildAmount, requestData.PetAmount)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, booking)
}

func (h Handler) AddMessage(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.BOOKING_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData AddMessageRequest
	err = helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")

	booking, err := h.application.AddMessage(identity.UserID, *recordID, requestData.Text)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, booking)
}
