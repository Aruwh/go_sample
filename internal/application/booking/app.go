package application

import (
	applicationAdminUser "fewoserv/internal/application/admin_user"
	applicationApartment "fewoserv/internal/application/apartment"
	applicationProcessLog "fewoserv/internal/application/process_log"
	applicationSaison "fewoserv/internal/application/saison"
	applicationSettings "fewoserv/internal/application/settings"
	applicationUser "fewoserv/internal/application/user"
	booking "fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	emailhandler "fewoserv/internal/interface/email_handler"
	repositoryBooking "fewoserv/internal/repository/booking"
	"fewoserv/pkg/mongodb"
	"fmt"
	"strconv"
	"time"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID, apartmentID string, message *string, status common.BookingtStatus, fromDate, toDate time.Time, adultAmount, childAmount, petAmount int) (*booking.Booking, error)
		Delete(userID, recordID string) error
		Get(recordID string) (*booking.Booking, error)
		GetBookingOverviews(userID *string, audianceType *common.AudianceType, date time.Time) ([]*booking.BookingOverview, error)
		GetBlockedDates(apartmentID string) ([]time.Time, error)

		Update(userID, recordID string, status *common.BookingtStatus, fromDate, toDate *time.Time, adultAmount, childAmount, petAmount *int) (*booking.Booking, error)

		CalculateBookingPrice(apartmentID string, fromDate, toDate time.Time) (*booking.PriceSummary, error)

		AddMessage(userID, recordID string, message string) (*booking.Booking, error)
	}

	Application struct {
		landingpageEndpoint *string
		copyEmailAddresses  []string
		repoBooking         *repositoryBooking.Repo
		processLog          applicationProcessLog.IApplication
		appSeason           applicationSaison.IApplication
		appSettings         applicationSettings.IApplication
		appApartment        applicationApartment.IApplication
		appUser             applicationUser.IApplication
		appAdminUser        applicationAdminUser.IApplication
		emailHandler        emailhandler.IEmailHandler
	}
)

func New(
	mongoDbClient mongodb.IClient,
	emailHandler emailhandler.IEmailHandler,
	processLog applicationProcessLog.IApplication,
	appSeason applicationSaison.IApplication,
	appSettings applicationSettings.IApplication,
	appApartment applicationApartment.IApplication,
	appUser applicationUser.IApplication,
	appAdminUser applicationAdminUser.IApplication,
	landingpageEndpoint *string,
	copyEmailAddresses []string) IApplication {
	application := Application{
		repoBooking:         repositoryBooking.New(mongoDbClient),
		emailHandler:        emailHandler,
		appSeason:           appSeason,
		appSettings:         appSettings,
		appApartment:        appApartment,
		appUser:             appUser,
		appAdminUser:        appAdminUser,
		processLog:          processLog,
		landingpageEndpoint: landingpageEndpoint,
		copyEmailAddresses:  copyEmailAddresses,
	}

	return &application
}

func (a *Application) isCreationAllowed(apartmentID string, bookingID *string, fromDate, toDate time.Time) error {
	usedFromDate, usedToDate := utils.SwapDatesIfNeeded(fromDate, toDate)

	bookingIDs, err := a.repoBooking.FindManyBetweenDate(apartmentID, bookingID, usedFromDate, usedToDate)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	if len(bookingIDs) > 0 {
		return ErrPlacingBookingNotAllowed
	}

	return nil
}

func (a *Application) sendEmailOnStatusChange(booking *booking.Booking, sendCopy bool) error {
	var (
		err       error
		template  email_template.IEmailTemplate
		usedEmail *string
		firstName *string
		lastName  *string
		locale    *common.Locale
	)

	user, err := a.appUser.Get(booking.UserID)
	if user != nil {
		usedEmail = user.Email
		firstName = user.FirstName
		lastName = user.LastName
		locale = user.Locale
	}

	adminUser, err := a.appAdminUser.Get(booking.UserID)
	if adminUser != nil {
		usedEmail = adminUser.Email
		firstName = adminUser.FirstName
		lastName = adminUser.LastName
		locale = adminUser.Locale
	}

	hasARealErrorOccoured := err != nil && usedEmail == nil
	if hasARealErrorOccoured {
		fmt.Errorf("%w: %v", ErrCantUpdate, err)
		return nil
	}

	apartment, err := a.appApartment.Get(booking.ApartmentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	// verfügbar > anfrage
	changed := booking.Status == common.Reserved
	if changed {
		template = BuildEmailIncomeResponseTemplate(*locale, *firstName, *lastName, *a.landingpageEndpoint, *apartment.Name, *booking)
	}

	// anfrage > anfragebestätigung
	changed = booking.Status == common.Confirmed
	if changed {
		template = BuildEmailRequestConfirmationTemplate(*locale, *firstName, *lastName, *a.landingpageEndpoint, *apartment.Name, *booking)
	}

	// x > storniert
	changed = booking.Status == common.Canceled
	if changed {
		template = BuildEmailIncomeCancelationTemplate(*locale, *firstName, *lastName, *a.landingpageEndpoint, *apartment.Name, *booking)
	}

	// it makes no sense to send an email without a valid template
	if template == nil {
		return nil
	}

	err = a.emailHandler.Send(*usedEmail, template)
	if err != nil {
		return err
	}

	// send copy to fewoserv
	if sendCopy {
		for _, emailAddress := range a.copyEmailAddresses {
			err = a.emailHandler.Send(emailAddress, template)
			if err != nil {
				return err
			}
		}

	}

	return err
}

func (a *Application) CalculateBookingPrice(apartmentID string, fromDate, toDate time.Time) (*booking.PriceSummary, error) {
	apartment, err := a.appApartment.Get(apartmentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	datesWithSeasonTypes, err := a.appSeason.LoadDatesWithSaisonTypes(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	newBooking := booking.CalculateBookingPrice("", common.Available, 0, apartment, fromDate, toDate, nil, datesWithSeasonTypes)

	return newBooking.PriceSummary, nil
}

func (a *Application) Create(userID, apartmentID string, message *string, status common.BookingtStatus, fromDate, toDate time.Time, adultAmount, childAmount, petAmount int) (*booking.Booking, error) {
	apartment, err := a.appApartment.Get(apartmentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.isCreationAllowed(apartmentID, nil, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	guestInfo := booking.NewGuestInfo(&adultAmount, &childAmount, &petAmount)

	datesWithSeasonTypes, err := a.appSeason.LoadDatesWithSaisonTypes(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	bookingNumber, err := a.appSettings.IncBookingNumber()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIncBookingNumber, err)
	}

	newBooking := booking.New(userID, status, bookingNumber, apartment, fromDate, toDate, guestInfo, datesWithSeasonTypes)
	if message != nil {
		newMessage := booking.Message{Timestamp: time.Now(), Text: *message}
		newBooking.Messages = append(newBooking.Messages, newMessage)
	}

	err = a.repoBooking.Insert(newBooking)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	value := fmt.Sprintf("booking number %s ", strconv.Itoa(newBooking.BookingNumber))
	a.processLog.New(userID, value, common.CREATED, common.BOOKING, &newBooking.ID)

	// send email
	err = a.sendEmailOnStatusChange(newBooking, true)
	if err != nil {
		fmt.Println("ERROR:", fmt.Errorf("%w: %v", ErrCantSave, err))
	}

	return newBooking, nil
}

func (a *Application) AddMessage(userID, recordID string, message string) (*booking.Booking, error) {
	foundBooking, err := a.repoBooking.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	newMessage := booking.Message{Timestamp: time.Now(), Text: message}
	foundBooking.Messages = append(foundBooking.Messages, newMessage)

	err = a.repoBooking.Update(userID, foundBooking)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	// TODO: send email
	a.processLog.New(userID, "message", common.CREATED, common.BOOKING, &foundBooking.ID)

	return foundBooking, nil
}

func (a *Application) Delete(userID, recordID string) error {
	booking, err := a.repoBooking.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repoBooking.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, strconv.Itoa(booking.BookingNumber), common.DELETED, common.BOOKING, &recordID)

	return nil
}

func (a *Application) Get(recordID string) (*booking.Booking, error) {
	foundBooking, err := a.repoBooking.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundBooking, nil
}

func (a *Application) GetMany(name string, sort common.Sort, skip, limit int64) ([]*booking.Booking, error) {
	foundBookings, err := a.repoBooking.FindMany(name, &sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundBookings, nil
}

func (a *Application) GetBookingOverviews(ownerID *string, audianceType *common.AudianceType, date time.Time) ([]*booking.BookingOverview, error) {
	apartments, err := a.appApartment.GetMany(ownerID, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	apartmentIDs := []*string{}
	for _, apartment := range apartments {
		apartmentIDs = append(apartmentIDs, &apartment.ID)
	}

	foundBookingsOverviews, err := a.repoBooking.FindManyBy(apartmentIDs, date)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	if foundBookingsOverviews == nil {
		foundBookingsOverviews = []*booking.BookingOverview{}
	}

	return foundBookingsOverviews, nil
}

func (a *Application) GetBlockedDates(apartmentID string) ([]time.Time, error) {
	year, month, _ := time.Now().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	endDate := time.Now().AddDate(1, 0, 0)

	return a.repoBooking.GetBlockedDates(apartmentID, startDate, endDate)
}

func (a *Application) Update(userID, recordID string, status *common.BookingtStatus, fromDate, toDate *time.Time, adultAmount, childAmount, petAmount *int) (*booking.Booking, error) {
	foundBooking, err := a.repoBooking.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.isCreationAllowed(foundBooking.ApartmentID, &recordID, *fromDate, *toDate)
	if err != nil {
		return nil, err
	}

	foundBooking.UpdateGuestInfo(adultAmount, childAmount, petAmount)
	hasStatusChanged, err := foundBooking.UpdateStatus(status)
	if err != nil {
		return nil, err
	}

	shouldDatesAreUpdated := fromDate != nil && toDate != nil
	if shouldDatesAreUpdated {
		refApartment, err := a.appApartment.Get(foundBooking.ApartmentID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
		}

		datesWithSeasonTypes, err := a.appSeason.LoadDatesWithSaisonTypes(*fromDate, *toDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
		}

		foundBooking.UpdateDates(*fromDate, *toDate, refApartment, datesWithSeasonTypes)
	}

	err = a.repoBooking.Update(userID, foundBooking)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, strconv.Itoa(foundBooking.BookingNumber), common.UPDATED, common.BOOKING, &recordID)

	if hasStatusChanged {
		err = a.sendEmailOnStatusChange(foundBooking, false)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
		}

		a.processLog.New(userID, strconv.Itoa(foundBooking.BookingNumber), common.SEND, common.BOOKING, &recordID)
	}

	return foundBooking, nil
}
