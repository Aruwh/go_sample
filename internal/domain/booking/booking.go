package booking

import (
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

type (
	Booking struct {
		ID            string `json:"id" bson:"_id"`
		BookingNumber int    `json:"bookingNumber" bson:"bookingNumber"`
		UserID        string `json:"userID" bson:"userID"`
		UserType      *common.AdminUserType
		ApartmentID   string                `json:"apartmentID" bson:"apartmentID"`
		Status        common.BookingtStatus `json:"status" bson:"status"`
		FromDate      time.Time             `json:"fromDate" bson:"fromDate"`
		ToDate        time.Time             `json:"toDate" bson:"toDate"`
		StayDays      int                   `json:"stayDays" bson:"stayDays"`
		GuestInfo     *GuestInfo            `json:"guestInfo" bson:"guestInfo"`
		PriceSummary  *PriceSummary         `json:"priceSummary" bson:"priceSummary"`
		Messages      []Message             `json:"messages" bson:"messages"`
		Created       shared.TimeStamp      `json:"created" bson:"created"`
		Edited        shared.TimeStamp      `json:"edited" bson:"edited"`
	}
)

func New(userID string, status common.BookingtStatus, bookingNumber int, apartment *apartment.Apartment, fromDate, toDate time.Time, guestInfo *GuestInfo, datesWithSeasonTypes map[time.Time]common.SaisonType) *Booking {
	timestamp := shared.NewTimeStamp(&userID)

	booking := Booking{
		ID:            mongodb.NewID(),
		BookingNumber: bookingNumber,
		UserID:        userID,
		ApartmentID:   apartment.ID,
		Status:        status,
		GuestInfo:     guestInfo,
		PriceSummary:  NewPriceSummary(),
		Messages:      []Message{},
		Edited:        timestamp,
		Created:       timestamp,
	}

	booking.calculateAndSet(fromDate, toDate, apartment, datesWithSeasonTypes)

	// we add one day, because the day of the depature counts as a half day
	booking.StayDays += 1

	return &booking
}

func CalculateBookingPrice(userID string, status common.BookingtStatus, bookingNumber int, apartment *apartment.Apartment, fromDate, toDate time.Time, guestInfo *GuestInfo, datesWithSeasonTypes map[time.Time]common.SaisonType) *Booking {
	timestamp := shared.NewTimeStamp(&userID)

	booking := Booking{
		ID:            mongodb.NewID(),
		BookingNumber: bookingNumber,
		UserID:        userID,
		ApartmentID:   apartment.ID,
		Status:        status,
		GuestInfo:     guestInfo,
		PriceSummary:  NewPriceSummary(),
		Messages:      []Message{},
		Edited:        timestamp,
		Created:       timestamp,
	}

	booking.calculateAndSet(fromDate, toDate, apartment, datesWithSeasonTypes)

	return &booking
}

func includes(arr []common.BookingtStatus, target common.BookingtStatus) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

func (b *Booking) calculateAndSet(fromDate, toDate time.Time, apartment *apartment.Apartment, datesWithSeasonTypes map[time.Time]common.SaisonType) {
	usedFromDate, usedToDate := utils.SwapDatesIfNeeded(fromDate, toDate)
	b.FromDate = usedFromDate
	b.ToDate = usedToDate

	b.StayDays = b.calcAndSetStayDays(usedFromDate, usedToDate)

	b.calcPriceSummary(apartment, datesWithSeasonTypes)
}

func (b *Booking) calcAndSetStayDays(fromDate, toDate time.Time) int {
	usedFromDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC)
	usedToDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.UTC)

	difference := usedToDate.Sub(usedFromDate)

	days := int(difference.Hours() / 24)

	return days

}

func (b *Booking) validateStatusSwitch(status *common.BookingtStatus) error {
	acceptableTransitions := map[common.BookingtStatus][]common.BookingtStatus{
		common.Available:      {common.Reserved, common.Blocked, common.BlockedByAdmin},
		common.Reserved:       {common.Confirmed, common.Canceled},
		common.Confirmed:      {common.Canceled, common.Completed},
		common.Blocked:        {common.Canceled},
		common.BlockedByAdmin: {common.Canceled},
		common.Canceled:       {},
		common.Completed:      {},
	}

	// Check if the current status allows the transition to the target status
	validTransitions, found := acceptableTransitions[b.Status]
	if !found {
		return fmt.Errorf("%w: invalid current status", ErrStatusSwitchNotAllowed)
	}

	isValid := includes(validTransitions, *status)
	if !isValid {
		return fmt.Errorf("%w: %d > %d", ErrStatusSwitchNotAllowed, b.Status, *status)
	}

	return nil
}

func (b *Booking) calcPriceSummary(apartment *apartment.Apartment, datesWithSeasonTypes map[time.Time]common.SaisonType) *PriceSummary {
	return b.PriceSummary.CalculatePriceSummary(b.FromDate, b.ToDate, apartment.SaisonPrice, datesWithSeasonTypes)
}

func (b *Booking) UpdateGuestInfo(adultAmount, childAmount, petAmount *int) {
	shouldBeUpdated := adultAmount != nil && *adultAmount != b.GuestInfo.AdultAmount
	if shouldBeUpdated {
		b.GuestInfo.AdultAmount = *adultAmount
	}

	shouldBeUpdated = childAmount != nil && *childAmount != b.GuestInfo.ChildAmount
	if shouldBeUpdated {
		b.GuestInfo.ChildAmount = *childAmount
	}

	shouldBeUpdated = petAmount != nil && *petAmount != b.GuestInfo.PetAmount
	if petAmount != nil {
		b.GuestInfo.PetAmount = *petAmount
	}
}

func (b *Booking) UpdateStatus(status *common.BookingtStatus) (bool, error) {
	hasStatusChanged := false

	shouldBeUpdated := status != nil && b.Status != *status
	if shouldBeUpdated {
		err := b.validateStatusSwitch(status)
		if err != nil {
			return hasStatusChanged, err
		}

		b.Status = *status
		hasStatusChanged = true
	}

	return hasStatusChanged, nil
}

func (b *Booking) UpdateDates(fromDate, toDate time.Time, refApartment *apartment.Apartment, datesWithSeasonTypes map[time.Time]common.SaisonType) error {
	hasSomethingChanged := false

	usedFromDate, usedToDate := utils.SwapDatesIfNeeded(fromDate, toDate)

	shouldBeUpdated := b.FromDate != usedFromDate
	if shouldBeUpdated {
		b.FromDate = fromDate
		hasSomethingChanged = true
	}

	shouldBeUpdated = b.ToDate != usedToDate
	if shouldBeUpdated {
		b.ToDate = toDate
		hasSomethingChanged = true
	}

	if hasSomethingChanged {
		b.calculateAndSet(fromDate, toDate, refApartment, datesWithSeasonTypes)
	}

	return nil
}
