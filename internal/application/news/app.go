package application

import (
	processLogApp "fewoserv/internal/application/process_log"
	"fmt"

	news "fewoserv/internal/domain/news"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/news"
	"fewoserv/pkg/mongodb"
	"time"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Create(userID string, title, content shared.Translation, publishAt time.Time) (*news.News, error)
		Delete(userID, recordID string) error
		Get(recordID string) (*news.News, error)
		GetMany(title *string, sort common.Sort, skip, limit int64) ([]*news.News, error)
		Update(userID, recordID string, title, content *shared.Translation, publishAt *time.Time, active *bool) (*news.News, error)

		GetManyPublic() ([]*news.News, error)
	}

	Application struct {
		processLog processLogApp.IApplication
		repo       *repository.Repo
	}
)

func New(mongoDbClient mongodb.IClient, processLog processLogApp.IApplication) IApplication {
	application := Application{repo: repository.New(mongoDbClient), processLog: processLog}

	return &application
}

func (a *Application) Create(userID string, name, content shared.Translation, publishAt time.Time) (*news.News, error) {
	currentDate := time.Now()
	isActive := currentDate.After(publishAt) || currentDate.Equal(publishAt)

	newApartment := news.New(userID, name, content, publishAt, isActive)

	err := a.repo.Insert(newApartment)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, *newApartment.Title.De_DE, common.CREATED, common.NEWS, &newApartment.ID)

	return newApartment, nil
}

func (a *Application) Delete(userID, recordID string) error {
	news, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, *news.Title.De_DE, common.DELETED, common.NEWS, &recordID)

	return nil
}

func (a *Application) Get(recordID string) (*news.News, error) {
	foundNews, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundNews, nil
}

func (a *Application) GetMany(title *string, sort common.Sort, skip, limit int64) ([]*news.News, error) {
	foundNews, err := a.repo.FindMany(title, &sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundNews, nil
}

func (a *Application) Update(userID, recordID string, name, content *shared.Translation, publishAt *time.Time, active *bool) (*news.News, error) {
	foundNews, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	foundNews.Update(name, content, publishAt, active)

	err = a.repo.Update(userID, foundNews)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, *foundNews.Title.De_DE, common.UPDATED, common.NEWS, &recordID)

	return foundNews, nil
}

func (a *Application) GetManyPublic() ([]*news.News, error) {
	foundNews, err := a.repo.FindManyPublic()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundNews, nil
}
