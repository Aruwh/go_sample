package public

import (
	application "fewoserv/internal/application/public"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strings"
	"time"
)

var (
	DOMAIN = "Public"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		GetApartments(w http.ResponseWriter, r *http.Request)
		GetApartmentDetails(w http.ResponseWriter, r *http.Request)
		PlaceBooking(w http.ResponseWriter, r *http.Request)
		CalculateBookingPrice(w http.ResponseWriter, r *http.Request)
		GetPictures(w http.ResponseWriter, r *http.Request)
		GetNews(w http.ResponseWriter, r *http.Request)
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

func (h Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	apartmentsOverview, err := h.application.GetNews()
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartmentsOverview)
}

func (h Handler) GetApartments(w http.ResponseWriter, r *http.Request) {
	apartmentsOverview, err := h.application.GetApartmentsOverview(common.MIDDLE)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartmentsOverview)
}

func (h Handler) GetApartmentDetails(w http.ResponseWriter, r *http.Request) {
	recordID := helper.GetQueryVar(r, "id")

	apartmentsOverview, err := h.application.GetApartmentDetails(*recordID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, apartmentsOverview)
}

func (h Handler) PlaceBooking(w http.ResponseWriter, r *http.Request) {
	var requestData PlaceBookingRequest
	err := helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	unixFromDate, err := utils.ConvertUnixStringToInt(&requestData.FromDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}
	fromDate := time.Unix(int64(unixFromDate), 0)

	unixToDate, err := utils.ConvertUnixStringToInt(&requestData.ToDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}
	toDate := time.Unix(int64(unixToDate), 0)

	unixBirthDate, err := utils.ConvertUnixStringToInt(&requestData.UserData.BirthDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}
	birthDate := time.Unix(int64(unixBirthDate), 0)

	placementID, err := h.application.PlaceBooking(
		requestData.ApartmentID,
		fromDate,
		toDate,
		birthDate,
		requestData.UserData.Sex,
		requestData.UserData.Locale,
		requestData.GuestInfo.AdultAmount,
		requestData.GuestInfo.ChildAmount,
		requestData.GuestInfo.PetAmount,
		requestData.UserID,
		requestData.Message,
		requestData.UserData.Address.Street,
		requestData.UserData.Address.StreetNumber,
		requestData.UserData.Address.Country,
		requestData.UserData.Address.City,
		requestData.UserData.PhoneNumber,
		requestData.UserData.Email,
		requestData.UserData.FirstName,
		requestData.UserData.LastName,
		requestData.UserData.Address.PostCode,
	)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, placementID)
}

func (h Handler) CalculateBookingPrice(w http.ResponseWriter, r *http.Request) {
	var requestData CalcBookingPriceRequest

	apartmentID := helper.GetQueryParam(r, "apartmentID")
	requestData.ApartmentID = *apartmentID

	urlFromDate := helper.GetQueryParam(r, "fromDate")
	unixFromDate, err := utils.ConvertUnixStringToInt(urlFromDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	urlToDate := helper.GetQueryParam(r, "toDate")
	unixToDate, err := utils.ConvertUnixStringToInt(urlToDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	fromDate := time.Unix(int64(unixFromDate), 0)
	toDate := time.Unix(int64(unixToDate), 0)

	priceSummary, err := h.application.CalculateBookingPrice(requestData.ApartmentID, fromDate, toDate)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, priceSummary)

}

func (h Handler) GetPictures(w http.ResponseWriter, r *http.Request) {
	var requestData GetPicturesRequest

	urlPictureIDs := helper.GetQueryParam(r, "pictureIDs")
	if urlPictureIDs == nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, nil, nil)
		return
	}
	requestData.PictureIDs = strings.Split(*urlPictureIDs, ",")

	variant := helper.GetQueryParam(r, "variant")
	if variant == nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, nil, nil)
		return
	}
	requestData.Variant = common.PictureVariant(*variant)

	// cachedPicture := h.pictureCache.Get(*variant, *recordID)
	// if cachedPicture != nil {
	// 	h.response.CreateSuccess(r.Context(), w, http.StatusOK, cachedPicture)
	// 	return
	// }

	picture, err := h.application.GetPictures(requestData.PictureIDs, requestData.Variant)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	// h.pictureCache.Add(*variant, picture)

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, picture)

}
