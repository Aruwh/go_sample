package application

import (
	appApartment "fewoserv/internal/application/apartment"
	appAttribute "fewoserv/internal/application/attribute"
	appBooking "fewoserv/internal/application/booking"
	appNews "fewoserv/internal/application/news"
	appPicture "fewoserv/internal/application/picture"
	appProcessLog "fewoserv/internal/application/process_log"
	appUser "fewoserv/internal/application/user"
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/domain/booking"
	"fewoserv/internal/domain/news"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/domain/user"
	"time"

	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fmt"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		GetApartmentDetails(apartmentID string) (*ApartmentDetail, error)
		GetApartmentsOverview(variant common.PictureVariant) ([]ApartmentOverview, error)

		PlaceBooking(apartmentID string, fromDate, toDate, birthDate time.Time, sex common.Sex, locale common.Locale, adultAmount, childAmoun, petAmount int, userID, message, street, streetNumber, country, city, phoneNumber, email, firstName, lastName, postCode *string) (*int, error)

		CalculateBookingPrice(apartmentID string, fromDate, toDate time.Time) (*booking.PriceSummary, error)
		GetPictures(pictureIDs []string, variant common.PictureVariant) ([]*shared.Picture, error)
		GetNews() ([]*news.News, error)
	}

	Application struct {
		processLog appProcessLog.IApplication
		attribute  appAttribute.IApplication
		picture    appPicture.IApplication
		apartment  appApartment.IApplication
		booking    appBooking.IApplication
		news       appNews.IApplication
		user       appUser.IApplication
	}

	// list
	ApartmentOverviewPicture struct {
		Source      string             `json:"source"`
		Description shared.Translation `json:"description"`
	}
	ApartmentOverviewAttribute struct {
		Svg  string             `json:"svg"`
		Name shared.Translation `json:"name"`
	}
	ApartmentOverviewAttributeInfo struct {
		Attributes            []ApartmentOverviewAttribute `json:"attributes"`
		AllowedNumberOfPeople int                          `json:"allowedNumberOfPeople"`
		AllowedNumberOfPets   int                          `json:"allowedNumberOfPets"`
		Bathrooms             int                          `json:"bathrooms"`
		RoomSize              int                          `json:"roomSize"`
		SleepingPlaces        int                          `json:"sleepingPlaces"`
	}
	ApartmentOverview struct {
		ApartmentID   string                            `json:"apartmentID"`
		Name          string                            `json:"name"`
		PictureIDs    []string                          `json:"pictureIDs"`
		AttributeInfo ApartmentOverviewAttributeInfo    `json:"attributeInfo"`
		SeasonPrices  map[common.SaisonType]interface{} `json:"seasonPrice"`
	}

	ApartmentDetail struct {
		ID            string                            `json:"id" bson:"_id"`
		Name          *string                           `json:"name" bson:"name"`
		Description   *shared.Translation               `json:"description" bson:"description"`
		PictureIDs    []string                          `json:"pictureIDs" bson:"pictureIDs"`
		AttributeInfo ApartmentOverviewAttributeInfo    `json:"attributeInfo"`
		SeasonPrices  map[common.SaisonType]interface{} `json:"seasonPrice"`
		BlockedDates  []time.Time                       `json:"blockedDates"`
	}
)

func New(apartment appApartment.IApplication, attribute appAttribute.IApplication, picture appPicture.IApplication, booking appBooking.IApplication, processLog appProcessLog.IApplication, news appNews.IApplication, user appUser.IApplication) IApplication {
	application := Application{
		processLog: processLog,
		apartment:  apartment,
		attribute:  attribute,
		picture:    picture,
		booking:    booking,
		news:       news,
		user:       user,
	}

	return &application
}

func (a *Application) getBlockedDates(apartmentID string) ([]time.Time, error) {
	return a.booking.GetBlockedDates(apartmentID)
}

func (a *Application) buildSeasonPrices(apartment *apartment.Apartment) map[common.SaisonType]interface{} {
	seasonPrices := make(map[common.SaisonType]interface{})

	seasonTypes := []common.SaisonType{common.SaisonLow, common.SaisonMiddle, common.SaisonHigh, common.SaisonPeak}
	for _, seasonType := range seasonTypes {
		bruttoPriceInfo := apartment.SaisonPrice.GetBruttoPriceInfo(seasonType, true)

		seasonPrices[seasonType] = bruttoPriceInfo.Price
	}

	return seasonPrices
}

func (a *Application) buildApartmentAttributeInfo(apartment *apartment.Apartment) (*ApartmentOverviewAttributeInfo, error) {
	attributeInfo := ApartmentOverviewAttributeInfo{}

	attributeInfo.AllowedNumberOfPeople = apartment.AttributeCollection.AllowedNumberOfPeople
	attributeInfo.AllowedNumberOfPets = apartment.AttributeCollection.AllowedNumberOfPets
	attributeInfo.RoomSize = apartment.AttributeCollection.RoomSize
	attributeInfo.SleepingPlaces = apartment.AttributeCollection.SleepingPlaces
	attributeInfo.Bathrooms = apartment.AttributeCollection.Bathrooms
	attributeInfo.Attributes = []ApartmentOverviewAttribute{}

	attributeIDs := append(apartment.AttributeCollection.AttributeIDs, apartment.AttributeCollection.TopAttributeIDs...)
	attributes, err := a.attribute.GetManyByID(attributeIDs)
	if err != nil {
		return nil, err
	}

	for _, attribute := range attributes {
		apartmentOverviewAttribute := ApartmentOverviewAttribute{}
		apartmentOverviewAttribute.Name = *attribute.Name

		apartmentOverviewAttribute.Svg = ""
		if attribute.Svg != nil {
			apartmentOverviewAttribute.Svg = *attribute.Svg
		}

		attributeInfo.Attributes = append(attributeInfo.Attributes, apartmentOverviewAttribute)
	}

	return &attributeInfo, nil
}

func (a *Application) GetApartmentsOverview(variant common.PictureVariant) ([]ApartmentOverview, error) {
	foundApartments, err := a.apartment.GetManyPublic()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	apartmentOverviews := []ApartmentOverview{}

	for _, apartment := range foundApartments {
		apartmentOverView := ApartmentOverview{
			ApartmentID:  apartment.ID,
			Name:         *apartment.Name,
			SeasonPrices: make(map[common.SaisonType]interface{}),
			AttributeInfo: ApartmentOverviewAttributeInfo{
				Attributes:            []ApartmentOverviewAttribute{},
				AllowedNumberOfPeople: 0,
				RoomSize:              0,
				SleepingPlaces:        0,
			},
		}

		apartmentOverView.PictureIDs = apartment.PictureIDs

		// attributes
		attributeInfo, err := a.buildApartmentAttributeInfo(apartment)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
		}
		apartmentOverView.AttributeInfo = *attributeInfo

		// seasonTypes
		apartmentOverView.SeasonPrices = a.buildSeasonPrices(apartment)

		apartmentOverviews = append(apartmentOverviews, apartmentOverView)
	}

	return apartmentOverviews, nil
}

func (a *Application) PlaceBooking(apartmentID string, fromDate, toDate, birthDate time.Time, sex common.Sex, locale common.Locale, adultAmount, childAmoun, petAmount int, userID, message, street, streetNumber, country, city, phoneNumber, email, firstName, lastName, postCode *string) (*int, error) {
	usedUserID := shared.DEFAULT_USER_ID
	if userID == nil {
		// TODO: load user
		address := a.user.CreateAddress(*street, *streetNumber, *country, *city, *firstName, *lastName, *postCode)
		newAddresses := []user.Address{address}

		phoneNumber := a.user.CreatePhoneNumber(common.Standard, *phoneNumber)
		newPhoneNumbers := []user.PhoneNumber{phoneNumber}

		newUser, err := a.user.Create(*email, sex, birthDate, locale, *firstName, *lastName, newAddresses, newPhoneNumbers)
		if err != nil {
			return nil, err
		}
		usedUserID = *&newUser.ID
	}

	booking, err := a.booking.Create(usedUserID, apartmentID, message, common.Reserved, fromDate, toDate, adultAmount, childAmoun, petAmount)
	if err != nil {
		return nil, err
	}

	// TODO send email

	return &booking.BookingNumber, nil
}

func (a *Application) CalculateBookingPrice(apartmentID string, fromDate, toDate time.Time) (*booking.PriceSummary, error) {
	priceSummary, err := a.booking.CalculateBookingPrice(apartmentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return priceSummary, nil
}

func (a *Application) GetPictures(pictureIDs []string, variant common.PictureVariant) ([]*shared.Picture, error) {
	return a.picture.GetMany(pictureIDs, variant)
}

func (a *Application) GetNews() ([]*news.News, error) {
	return a.news.GetManyPublic()
}

func (a *Application) GetApartmentDetails(apartmentID string) (*ApartmentDetail, error) {
	apartment, err := a.apartment.Get(apartmentID)
	if err != nil {
		return nil, err
	}

	// attributes
	attributeInfo, err := a.buildApartmentAttributeInfo(apartment)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	// blocked dates
	blockedDates, err := a.getBlockedDates(apartmentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	apartmentDetails := ApartmentDetail{
		ID:            apartment.ID,
		Name:          apartment.Name,
		Description:   apartment.Description,
		PictureIDs:    apartment.PictureIDs,
		AttributeInfo: *attributeInfo,
		SeasonPrices:  a.buildSeasonPrices(apartment),
		BlockedDates:  blockedDates,
	}

	return &apartmentDetails, nil
}
