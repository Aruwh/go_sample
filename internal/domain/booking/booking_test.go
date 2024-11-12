package booking_test

import (
	"fmt"
	"testing"
)

// import (
// 	"fewoserv/internal/domain/apartment"
// 	"fewoserv/internal/domain/booking"
// 	"fewoserv/internal/infrastructure/common"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	testApartment     = buildTestApartment()
// 	userID            = "asdadas"
// 	saison            = common.SaisonLow
// 	fromDate          = time.Date(2023, time.October, 15, 0, 0, 0, 0, time.UTC)
// 	toDate            = time.Date(2023, time.October, 21, 0, 0, 0, 0, time.UTC)
// 	gestInfo          = buildGuestInfo()
// 	message           = "grüezi mitenand"
// 	adultAmount   int = 6
// 	childAmount   int = 2
// 	petAmount     int = 1
// )

// func buildTestApartment() *apartment.Apartment {
// 	testApartment := apartment.New("asjdlasd", "asjdlasd", "3424jasda", "TestApartment")

// 	return testApartment
// }

// func buildGuestInfo() *booking.GuestInfo {
// 	return booking.NewGuestInfo(&adultAmount, &childAmount, &petAmount)
// }

// func TestNew(t *testing.T) {
// 	testBooking := booking.New(userID, testApartment, saison, fromDate, toDate, gestInfo, &message)

// 	assert.NotNil(t, testBooking)
// 	assert.Equal(t, testBooking.FromDate, fromDate)
// 	assert.Equal(t, testBooking.ToDate, toDate)
// 	assert.Equal(t, testBooking.StayDays, int(6))
// 	assert.NotNil(t, testBooking.GuestInfo)
// 	assert.Equal(t, testBooking.GuestInfo.AdultAmount, int(6))
// 	assert.Equal(t, testBooking.GuestInfo.ChildAmount, int(2))
// 	assert.Equal(t, testBooking.GuestInfo.PetAmount, int(1))
// 	assert.Equal(t, testBooking.Status, common.Available)
// 	assert.Equal(t, testBooking.Message, message)
// 	assert.NotNil(t, testBooking.PriceSummary)
// 	assert.Equal(t, testBooking.PriceSummary.Vat, int(7))
// 	assert.Equal(t, testBooking.PriceSummary.Total, float32(963))
// 	assert.Equal(t, testBooking.PriceSummary.Tax, float32(63))
// }

// func TestUpdateGuestInfo(t *testing.T) {
// 	testBooking := booking.New(userID, testApartment, saison, fromDate, toDate, gestInfo, &message)

// 	var newPetAmount int = 3
// 	testBooking.UpdateGuestInfo(nil, nil, &newPetAmount)

// 	assert.Equal(t, testBooking.GuestInfo.AdultAmount, adultAmount)
// 	assert.Equal(t, testBooking.GuestInfo.ChildAmount, childAmount)
// 	assert.Equal(t, testBooking.GuestInfo.PetAmount, newPetAmount)
// }

func TestUpdateGuestInfo(t *testing.T) {
	blub := make(map[int]string)
	blub[2] = "B"
	blub[3] = "C"
	blub[4] = "D"


	c, isOk := blub[2]
	fmt.Println(c)
	fmt.Println(isOk)

	c, isOk = blub[1]
	fmt.Println(c)
	fmt.Println(isOk)

}
