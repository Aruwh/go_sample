package application

import (
	appProcessLog "fewoserv/internal/application/process_log"
	apartment "fewoserv/internal/domain/apartment"

	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/apartment"
	"fewoserv/pkg/mongodb"
	"fmt"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID, ownerID, realEstateID, name string, description *shared.Translation, pictureIDs, attributeIDs, topAttributeIDs *[]string, saisonPrice *apartment.SaisonPrice, roomSize, sleepingPlaces, bathRooms, allowedNumberOfPeople, allowedNumberOfPets int) (*apartment.Apartment, error)
		Delete(userID, recordID string) error
		Get(recordID string) (*apartment.Apartment, error)
		GetReadOnly(ownerID, recordID string) (*apartment.ReadOnlyApartment, error)
		GetMany(ownerID, realEstateID, name *string, sort *common.Sort, skip, limit *int64) ([]*apartment.MinimalApartment, error)
		Update(userID, recordID string, isActive *bool, ownerID, realEstateID, name *string, description *shared.Translation, pictureIDs, attributeIDs, topAttributeIDs *[]string, saisonPrice *apartment.SaisonPrice, roomSize, sleepingPlaces, bathRooms, allowedNumberOfPeople, allowedNumberOfPets *int) (*apartment.Apartment, []string, error)

		GetManyPublic() ([]*apartment.Apartment, error)
	}

	Application struct {
		processLog appProcessLog.IApplication
		repo       *repository.Repo
	}
)

func New(mongoDbClient mongodb.IClient, processLog appProcessLog.IApplication) IApplication {
	application := Application{repo: repository.New(mongoDbClient), processLog: processLog}
	return &application
}

func (a *Application) ensureSwitchToActiveIsAllowed(apartment apartment.Apartment) error {
	areMinimumOfPicturesAttached := len(apartment.PictureIDs) >= 6
	if !areMinimumOfPicturesAttached {
		return ErrNotEnaughPictures
	}

	return nil
}

func (a *Application) Create(userID, ownerID, realEstateID, name string, description *shared.Translation, pictureIDs, attributeIDs, topAttributeIDs *[]string, saisonPrice *apartment.SaisonPrice, roomSize, sleepingPlaces, bathRooms, allowedNumberOfPeople, allowedNumberOfPets int) (*apartment.Apartment, error) {
	newApartment := apartment.New(userID, ownerID, realEstateID, name)

	newApartment.UpdateDescription(description)
	newApartment.UpdatePictureIDs(pictureIDs)
	newApartment.UpdateSaisonPrice(saisonPrice)
	newApartment.UpdateAttributeIDs(attributeIDs)
	newApartment.UpdateTopAttributeIDs(topAttributeIDs)
	newApartment.AttributeCollection.AllowedNumberOfPeople = allowedNumberOfPeople
	newApartment.AttributeCollection.AllowedNumberOfPets = allowedNumberOfPets
	newApartment.AttributeCollection.RoomSize = roomSize
	newApartment.AttributeCollection.SleepingPlaces = sleepingPlaces
	newApartment.AttributeCollection.Bathrooms = bathRooms

	err := a.repo.Insert(newApartment)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, name, common.CREATED, common.APARTMENT, &newApartment.ID)

	return newApartment, nil
}

func (a *Application) Delete(userID, recordID string) error {
	apartment, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v: %s", ErrCantDelete, err, recordID)
	}

	a.processLog.New(userID, *apartment.Name, common.DELETED, common.APARTMENT, &recordID)

	return nil
}

func (a *Application) Get(recordID string) (*apartment.Apartment, error) {
	foundApartment, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundApartment, nil
}

func (a *Application) GetReadOnly(ownerID, recordID string) (*apartment.ReadOnlyApartment, error) {
	foundApartment, err := a.repo.GetReadOnly(ownerID, recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundApartment, nil
}

func (a *Application) GetMany(ownerID, realEstateID, name *string, sort *common.Sort, skip, limit *int64) ([]*apartment.MinimalApartment, error) {
	foundApartments, err := a.repo.FindManyMinimal(ownerID, realEstateID, name, sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundApartments, nil
}

func (a *Application) Update(userID, recordID string, isActive *bool, ownerID, realEstateID, name *string, description *shared.Translation, pictureIDs, attributeIDs, topAttributeIDs *[]string, saisonPrice *apartment.SaisonPrice, roomSize, sleepingPlaces, bathRooms, allowedNumberOfPeople, allowedNumberOfPets *int) (*apartment.Apartment, []string, error) {
	foundApartment, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, []string{}, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	foundApartment.UpdateName(name)
	foundApartment.UpdateDescription(description)
	foundApartment.UpdateSaisonPrice(saisonPrice)
	foundApartment.UpdateAttributeIDs(attributeIDs)
	foundApartment.UpdateTopAttributeIDs(topAttributeIDs)

	pictureIDsToRemove := foundApartment.UpdatePictureIDs(pictureIDs)

	shouldBeUpdated := isActive != nil && foundApartment.IsActive != *isActive
	if shouldBeUpdated {
		err := a.ensureSwitchToActiveIsAllowed(*foundApartment)
		if err != nil {
			return nil, []string{}, err
		}

		foundApartment.IsActive = *isActive
	}

	shouldBeUpdated = ownerID != nil && *foundApartment.OwnerID != *ownerID
	if shouldBeUpdated {
		foundApartment.OwnerID = ownerID
	}

	shouldBeUpdated = realEstateID != nil && *foundApartment.RealEstateID != *realEstateID
	if shouldBeUpdated {
		foundApartment.RealEstateID = realEstateID
	}

	shouldBeUpdated = allowedNumberOfPeople != nil && foundApartment.AttributeCollection.AllowedNumberOfPeople != *allowedNumberOfPeople
	if shouldBeUpdated {
		foundApartment.AttributeCollection.AllowedNumberOfPeople = *allowedNumberOfPeople
	}

	shouldBeUpdated = allowedNumberOfPets != nil && foundApartment.AttributeCollection.AllowedNumberOfPets != *allowedNumberOfPets
	if shouldBeUpdated {
		foundApartment.AttributeCollection.AllowedNumberOfPets = *allowedNumberOfPets
	}

	shouldBeUpdated = roomSize != nil && foundApartment.AttributeCollection.RoomSize != *roomSize
	if shouldBeUpdated {
		foundApartment.AttributeCollection.RoomSize = *roomSize
	}

	shouldBeUpdated = sleepingPlaces != nil && foundApartment.AttributeCollection.SleepingPlaces != *sleepingPlaces
	if shouldBeUpdated {
		foundApartment.AttributeCollection.SleepingPlaces = *sleepingPlaces
	}

	shouldBeUpdated = bathRooms != nil && foundApartment.AttributeCollection.Bathrooms != *bathRooms
	if shouldBeUpdated {
		foundApartment.AttributeCollection.Bathrooms = *bathRooms
	}

	err = a.repo.Update(userID, foundApartment)
	if err != nil {
		return nil, []string{}, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, "", common.UPDATED, common.APARTMENT, &recordID)

	return foundApartment, pictureIDsToRemove, nil
}

func (a *Application) GetManyPublic() ([]*apartment.Apartment, error) {
	foundApartments, err := a.repo.FindManyPublic(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundApartments, nil
}
