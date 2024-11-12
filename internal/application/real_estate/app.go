package application

import (
	apartmentApp "fewoserv/internal/application/apartment"
	processLogApp "fewoserv/internal/application/process_log"
	realEstate "fewoserv/internal/domain/real_estate"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/real_estate"
	"fewoserv/pkg/mongodb"
	"fmt"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID, name string, pictureID *string, description shared.Translation) (*realEstate.RealEstate, error)
		Delete(userID, recordID string) error
		Get(recordID string) (*realEstate.RealEstate, error)
		GetMany(name *string, sort common.Sort, skip, limit int64) ([]*realEstate.RealEstate, error)
		Update(userID, recordID string, pictureID, name *string, description *shared.Translation) (*realEstate.RealEstate, *string, error)
	}

	Application struct {
		repo                 *repository.Repo
		applicationApartment apartmentApp.IApplication
		processLog           processLogApp.IApplication
	}
)

func New(mongoDbClient mongodb.IClient, processLog processLogApp.IApplication) IApplication {
	application := Application{
		repo:                 repository.New(mongoDbClient),
		applicationApartment: apartmentApp.New(mongoDbClient, processLog),
		processLog:           processLog,
	}

	return &application
}

func (a *Application) Create(userID, name string, pictureID *string, description shared.Translation) (*realEstate.RealEstate, error) {
	newRealEstate, err := realEstate.New(userID, &name, pictureID, &description)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantInitRealestate, err)
	}

	err = a.repo.Insert(newRealEstate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, name, common.CREATED, common.REAL_ESTATE, &newRealEstate.ID)

	return newRealEstate, nil
}

func (a *Application) Delete(userID, recordID string) error {
	realEstate, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	apartments, err := a.applicationApartment.GetMany(nil, &recordID, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	hasLinkedApartments := len(apartments) != 0
	if hasLinkedApartments {
		return ErrHasLinkedApartments
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, *realEstate.Name, common.DELETED, common.REAL_ESTATE, &recordID)
	return nil
}

func (a *Application) Get(recordID string) (*realEstate.RealEstate, error) {
	foundRealEstate, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundRealEstate, nil
}

func (a *Application) GetMany(name *string, sort common.Sort, skip, limit int64) ([]*realEstate.RealEstate, error) {
	foundRecords, err := a.repo.FindMany(nil, name, &sort, &skip, &limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundRecords, nil
}

func (a *Application) Update(userID, recordID string, pictureID, name *string, description *shared.Translation) (*realEstate.RealEstate, *string, error) {
	foundRealEstate, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	pictureIDToRemove := foundRealEstate.Update(pictureID, name, description)
	err = a.repo.Update(userID, foundRealEstate)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, *foundRealEstate.Name, common.UPDATED, common.REAL_ESTATE, &recordID)

	return foundRealEstate, pictureIDToRemove, nil
}
