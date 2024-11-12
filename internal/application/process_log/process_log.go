package application

import (
	process_log "fewoserv/internal/domain/process_log"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/process_log"
	"fewoserv/pkg/mongodb"
	"fmt"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		New(userID, value string, action common.Action, domain common.Domain, recordID *string) (*process_log.ProcessLog, error)
		GetManyBy(userID string) ([]*process_log.ProcessLog, error)
	}

	Application struct {
		repo *repository.Repo
	}
)

func New(mongoDbClient mongodb.IClient) IApplication {
	application := Application{repo: repository.New(mongoDbClient)}

	return &application
}

func (a *Application) New(createrID, value string, action common.Action, domain common.Domain, recordID *string) (*process_log.ProcessLog, error) {
	newProcessLog := process_log.New(createrID, value, action, domain, recordID)

	err := a.repo.Insert(newProcessLog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	return newProcessLog, nil
}

func (a *Application) GetManyBy(userID string) ([]*process_log.ProcessLog, error) {
	processLogs, err := a.repo.FindMany(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v: %s", ErrRecordNotExists, err, userID)
	}

	return processLogs, nil
}
