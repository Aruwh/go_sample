package application

import (
	booking "fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	firstName     = "Horst"
	lastName      = "Schlemmer"
	feDestination = "localhost:4000"
	apartmentName = "Dream Heaven"

	locale = common.DeDE

	testBooking = booking.Booking{
		ApartmentID:   "2354251345345345345",
		BookingNumber: 54,
		FromDate:      time.Now(),
		ToDate:        time.Now(),
	}
)

func TestBuildEmailIncomeResponseTemplate(t *testing.T) {
	template := BuildEmailIncomeResponseTemplate(locale, firstName, lastName, feDestination, apartmentName, testBooking)

	assert.NotNil(t, template)
}

func TestBuildEmailIncomeCancelationTemplate(t *testing.T) {
	template := BuildEmailIncomeCancelationTemplate(locale, firstName, lastName, feDestination, apartmentName, testBooking)

	assert.NotNil(t, template)
}

func TestBuildEmailRequestConfirmationTemplate(t *testing.T) {
	template := BuildEmailRequestConfirmationTemplate(locale, firstName, lastName, feDestination, apartmentName, testBooking)

	assert.NotNil(t, template)
}
