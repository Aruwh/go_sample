package application

import (
	"encoding/json"
	applicationProcessLog "fewoserv/internal/application/process_log"
	attribute "fewoserv/internal/domain/attribute"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/attribute"
	"fewoserv/pkg/mongodb"
	"fmt"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID string, name *shared.Translation, svg *string) (*attribute.Attribute, error)
		Delete(userID, recordID string) error
		Get(userID string) (*attribute.Attribute, error)
		GetMany(name *string, sort common.Sort, skip, limit int64) ([]*attribute.Attribute, error)
		Update(userID, recordID string, name *shared.Translation, svg *string) (*attribute.Attribute, error)
		CreateDefaults(userID string) error
		GetManyByID(recordIDs []string) ([]*attribute.Attribute, error)
	}

	Application struct {
		processLog applicationProcessLog.IApplication
		repo       *repository.Repo
	}
)

func New(mongoDbClient mongodb.IClient, processLog applicationProcessLog.IApplication) IApplication {
	application := Application{repo: repository.New(mongoDbClient), processLog: processLog}

	return &application
}

func (a *Application) CreateDefaults(userID string) error {
	// Define a struct to represent the structure of your JSON
	var translations map[string]shared.Translation

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(defaultTranslations), &translations)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantUnmarshal, err)
	}

	// Iterate over the JSON object
	for _, translationName := range translations {
		newAttribute := attribute.New(userID, &translationName, nil)

		entryAlreadyExists := a.repo.ValidateRecordExists(newAttribute)
		if entryAlreadyExists {
			continue
		}

		newAttribute, err := a.Create(userID, &translationName, nil)
		if err != nil {
			return err
		}

		a.processLog.New(userID, *newAttribute.Name.De_DE, common.CREATED, common.ATTRIBUTE, &newAttribute.ID)

	}

	return nil
}

func (a *Application) Create(userID string, name *shared.Translation, svg *string) (*attribute.Attribute, error) {
	newApartment := attribute.New(userID, name, svg)

	err := a.repo.Insert(newApartment)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, *name.De_DE, common.CREATED, common.ATTRIBUTE, &newApartment.ID)

	return newApartment, nil
}

func (a *Application) Delete(userID, recordID string) error {
	attribute, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, *attribute.Name.De_DE, common.DELETED, common.ATTRIBUTE, &recordID)

	return nil
}

func (a *Application) Get(userID string) (*attribute.Attribute, error) {
	foundAttribute, err := a.repo.LoadByID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundAttribute, nil
}

func (a *Application) GetMany(name *string, sort common.Sort, skip, limit int64) ([]*attribute.Attribute, error) {
	foundAttributes, err := a.repo.FindMany(name, &sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundAttributes, nil
}

func (a *Application) Update(userID, recordID string, name *shared.Translation, svg *string) (*attribute.Attribute, error) {
	foundAttribute, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	foundAttribute.Update(name, svg)

	err = a.repo.Update(recordID, foundAttribute)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, *foundAttribute.Name.De_DE, common.UPDATED, common.ATTRIBUTE, &recordID)

	return foundAttribute, nil
}

func (a *Application) GetManyByID(recordIDs []string) ([]*attribute.Attribute, error) {
	foundAttributes, err := a.repo.FindManyByID(recordIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundAttributes, nil
}
