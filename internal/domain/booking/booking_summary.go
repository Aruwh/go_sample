package booking

import (
	"time"
)

type (
	BookingMessage struct {
		Timestamp    time.Time `bson:"timestamp" json:"timestamp"`
		SendFromUser bool      `bson:"sendFromUser" json:"sendFromUser"`
		Text         string    `bson:"text" json:"text"`
	}

	BookingSummary struct {
		BookingID     string             `bson:"bookingID" json:"bookingID"`
		BookingNumber int                `bson:"bookingNumber" json:"bookingNumber"`
		Status        int                `bson:"status" json:"status"`
		FromDate      time.Time          `bson:"fromDate" json:"fromDate"`
		ToDate        time.Time          `bson:"toDate" json:"toDate"`
		StayDays      int                `bson:"stayDays" json:"stayDays"`
		UserID        string             `bson:"userID" json:"userID"`
		UserName      string             `bson:"userName" json:"userName"`
		GuestInfo     map[string]int     `bson:"guestInfo" json:"guestInfo"`
		PriceSummary  map[string]float64 `bson:"priceSummary" json:"priceSummary"`
		Messages      []BookingMessage   `bson:"messages" json:"messages"`
	}

	BookingOverview struct {
		ApartmentID string           `bson:"apartmentID" json:"apartmentID"`
		Summaries   []BookingSummary `bson:"summaries" json:"summaries"`
	}
)
