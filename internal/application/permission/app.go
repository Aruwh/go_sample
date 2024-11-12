package application

import (
	processLogApp "fewoserv/internal/application/process_log"

	"fewoserv/internal/domain/permission"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/permission"
	"fewoserv/pkg/mongodb"
	"fmt"
)

var (
	log                       = logger.New("APPLICATION")
	defaultRequestPermissions = []common.RequestPermission{
		common.APARTMENT_VIEW,
		common.APARTMENT_EDIT,
		common.APARTMENT_DELETE,
		common.BOOKING_VIEW,
		common.BOOKING_EDIT,
		common.BOOKING_DELETE,
		common.PERMISSION_VIEW,
		common.PERMISSION_EDIT,
		common.PERMISSION_DELETE,
		common.REAL_ESTATE_VIEW,
		common.REAL_ESTATE_EDIT,
		common.REAL_ESTATE_DELETE,
		common.NEWS_VIEW,
		common.NEWS_EDIT,
		common.NEWS_DELETE,
		common.SETTINGS_VIEW,
		common.SETTINGS_EDIT,
		common.SETTINGS_DELETE,
		common.USER_VIEW,
		common.USER_EDIT,
		common.USER_DELETE,
		common.PICTURE_VIEW,
		common.PICTURE_EDIT,
		common.PICTURE_DELETE,
	}
)

type (
	IApplication interface {
		CreateDefaults() []common.RequestPermission
		GetMany(name *string, sort common.Sort, skip, limit int64) ([]*permission.Permission, error)
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

func (a *Application) CreateDefaults() []common.RequestPermission {
	permissions := []common.RequestPermission{}

	for _, requestPermission := range defaultRequestPermissions {
		description := shared.NewTranslation("")

		newPermission := permission.New(string(requestPermission), *description)

		alreadyExists := a.repo.ValidatePermissionExists(newPermission)
		if alreadyExists {
			continue
		}

		err := a.repo.Insert(newPermission)
		if err == nil {
			log.Info(fmt.Sprintf("CreatPermissions :: %v %s created", requestPermission, newPermission.ID))
			permissions = append(permissions, requestPermission)
		} else {
			log.Error(fmt.Sprintf("CreatPermissions :: %v not created : %s", requestPermission, err))
		}
	}

	return permissions
}

func (a *Application) GetMany(name *string, sort common.Sort, skip, limit int64) ([]*permission.Permission, error) {
	foundApartments, err := a.repo.FindMany(name, &sort, skip, limit)
	if err != nil {
		fmt.Println(err)
		return nil, ErrRecordNotExists
	}

	return foundApartments, nil
}
