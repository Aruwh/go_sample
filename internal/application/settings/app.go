package application

import (
	processLogApp "fewoserv/internal/application/process_log"
	"fmt"

	"fewoserv/internal/domain/settings"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/settings"
	"fewoserv/pkg/mongodb"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		CreateDefaults() error
		GetBookingNumber() (int, error)
		IncBookingNumber() (int, error)
		GetNotificationMessage() (*shared.Translation, error)
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

func (a *Application) CreateDefaults() error {
	defaultSettings := settings.New()

	err := a.repo.Insert(defaultSettings)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrCantSave, err)
	}

	return nil
}

func (a *Application) GetBookingNumber() (int, error) {
	return a.repo.GetBookingNumber()
}

func (a *Application) IncBookingNumber() (int, error) {
	err := a.repo.IncBookingNumber()
	if err != nil {
		return 0, fmt.Errorf("%w, %v", ErrCantIncBookingNumber, err)
	}

	bookingNumber, err := a.GetBookingNumber()
	if err != nil {
		return 0, fmt.Errorf("%w, %v", ErrRecordNotExists, err)
	}

	return bookingNumber, nil
}

func (a *Application) GetNotificationMessage() (*shared.Translation, error) {
	message, err := a.repo.GetNotificationMessage()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return message, nil
}
