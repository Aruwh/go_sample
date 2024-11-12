package application

import (
	processLogApp "fewoserv/internal/application/process_log"
	"strconv"

	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/saison"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID string, year int, entries []shared.SaisonEntry) (*shared.Saison, error)
		Delete(userID, recordID string) error
		Get(recordID string) (*shared.Saison, error)
		GetMany(sort common.Sort, skip, limit int64, searchYear *int) ([]*shared.Saison, error)
		Update(userID, recordID string, year *int, entries *[]shared.SaisonEntry) (*shared.Saison, error)
		LoadDatesWithSaisonTypes(fromDate, toDate time.Time) (map[time.Time]common.SaisonType, error)
	}

	Application struct {
		processLog processLogApp.IApplication
		repo       *repository.Repo
	}
)

func (a *Application) validate(year int) error {
	foundSaisons, err := a.repo.FindMany(nil, 0, 0, &year)
	if err != nil {
		return fmt.Errorf("%w: %v", fmt.Errorf("%w: %v", ErrRecordNotExists, err), err)
	}

	if len(foundSaisons) > 0 {
		return fmt.Errorf("%w %v", ErrSeasonAlreadyExists, year)
	}

	return nil
}

func New(mongoDbClient mongodb.IClient, processLog processLogApp.IApplication) IApplication {
	application := Application{repo: repository.New(mongoDbClient), processLog: processLog}

	return &application
}

func (a *Application) Create(userID string, year int, entries []shared.SaisonEntry) (*shared.Saison, error) {
	err := a.validate(year)
	if err != nil {
		return nil, err
	}

	saison := shared.NewSaison(userID, year, entries)

	err = a.repo.Insert(saison)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, strconv.Itoa(year), common.CREATED, common.SEASON, &saison.ID)

	return saison, nil
}

func (a *Application) Delete(userID, recordID string) error {
	season, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, strconv.Itoa(season.Year), common.DELETED, common.SEASON, &recordID)
	return nil
}

func (a *Application) Get(recordID string) (*shared.Saison, error) {
	saison, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return saison, nil
}

func (a *Application) GetMany(sort common.Sort, skip, limit int64, searchYear *int) ([]*shared.Saison, error) {
	foundRecords, err := a.repo.FindMany(&sort, skip, limit, searchYear)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundRecords, nil
}

func (a *Application) LoadDatesWithSaisonTypes(fromDate, toDate time.Time) (map[time.Time]common.SaisonType, error) {
	years := []int{fromDate.Year(), toDate.Year()}

	foundSaison, err := a.repo.FindManyByYears(years)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	dateValueMap := make(map[time.Time]common.SaisonType)

	// pre fill map with values of the given years
	for _, year := range years {
		currentDate := time.Date(year, time.January, 1, 6, 0, 0, 0, time.UTC)

		for currentDate.Year() == year {
			dateValueMap[currentDate] = common.SaisonBase

			// jump to the next day
			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	// fill map corresponding of the season entries
	for _, saison := range foundSaison {
		for _, saisonEntry := range saison.Entries {
			currentDate := time.Date(saisonEntry.FromDate.Year(), saisonEntry.FromDate.Month(), saisonEntry.FromDate.Day(), 6, 0, 0, 0, time.UTC)

			for !currentDate.After(saisonEntry.ToDate) {
				dateValueMap[currentDate] = saisonEntry.Type

				currentDate = currentDate.AddDate(0, 0, 1)
			}
		}
	}

	return dateValueMap, nil
}

func (a *Application) Update(userID, recordID string, year *int, entries *[]shared.SaisonEntry) (*shared.Saison, error) {
	saison, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	saison.Edited.Update(userID)

	saison.Update(year, entries)
	err = a.repo.Update(userID, saison)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, strconv.Itoa(saison.Year), common.UPDATED, common.SEASON, &recordID)

	return saison, nil
}
