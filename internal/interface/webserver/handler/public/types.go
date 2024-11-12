package public

import (
	"fewoserv/internal/infrastructure/common"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	PlaceBookingGuestInfo struct {
		AdultAmount int `json:"adultAmount"`
		ChildAmount int `json:"childAmount"`
		PetAmount   int `json:"petAmount"`
	}

	PlaceBookingUserDataAddress struct {
		Street       *string `json:"street"`
		StreetNumber *string `json:"streetNumber"`
		Country      *string `json:"country"`
		City         *string `json:"city"`
		FirstName    *string `json:"firstName"`
		LastName     *string `json:"lastName"`
		PostCode     *string `json:"postCode"`
	}

	PlaceBookingUserData struct {
		FirstName   *string                     `json:"firstName"`
		LastName    *string                     `json:"lastName"`
		PhoneNumber *string                     `json:"phoneNumber"`
		Email       *string                     `json:"email"`
		Sex         common.Sex                  `json:"sex"`
		Locale      common.Locale               `json:"locale"`
		BirthDate   string                      `json:"birthDate"`
		Address     PlaceBookingUserDataAddress `json:"address"`
	}

	PlaceBookingRequest struct {
		UserID      *string               `json:"userID"`
		ApartmentID string                `json:"apartmentID" validate:"mongoDbID"`
		FromDate    string                `json:"fromDate"`
		ToDate      string                `json:"toDate"`
		Message     *string               `json:"message"`
		UserData    PlaceBookingUserData  `json:"userData"`
		GuestInfo   PlaceBookingGuestInfo `json:"guestInfo"`
	}

	CalcBookingPriceRequest struct {
		ApartmentID string `json:"apartmentID" validate:"mongoDbID"`
		FromDate    string `json:"fromDate"`
		ToDate      string `json:"toDate"`
	}

	GetPicturesRequest struct {
		PictureIDs []string              `json:"pictureIDs"`
		Variant    common.PictureVariant `json:"variant"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

)
