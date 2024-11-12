package booking

import (
	"fewoserv/internal/infrastructure/common"
	"time"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		ApartmentID string                `json:"apartmentID"`
		Status      common.BookingtStatus `json:"status"`
		FromDate    time.Time             `json:"fromDate"`
		ToDate      time.Time             `json:"toDate"`
		AdultAmount int                   `json:"adultAmount"`
		ChildAmount int                   `json:"childAmount"`
		PetAmount   int                   `json:"petAmount"`
	}

	UpdateRequest struct {
		Status      *common.BookingtStatus `json:"status"`
		AdultAmount *int                   `json:"adultAmount"`
		ChildAmount *int                   `json:"childAmount"`
		PetAmount   *int                   `json:"petAmount"`
		FromDate    *time.Time             `json:"fromDate"`
		ToDate      *time.Time             `json:"toDate"`
	}

	AddMessageRequest struct {
		Text string `json:"text"`
	}

	GetManyBookingOverviewsRequest struct {
		Date time.Time `json:"date"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
)
