package application

import (
	applicationProcessLog "fewoserv/internal/application/process_log"
	user "fewoserv/internal/domain/user"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	emailhandler "fewoserv/internal/interface/email_handler"
	repository "fewoserv/internal/repository/user"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

var (
	log = logger.New("APPLICATION")
)

type (
	PasswordUpdate struct {
		OldPassword        string
		NewPassword        string
		NewComparePassword string
	}

	IApplication interface {
		Create(
			email string,
			sex common.Sex,
			birthDate time.Time,
			locale common.Locale,
			firstName string,
			lastName string,
			addresses []user.Address,
			phoneNumbers []user.PhoneNumber,
		) (*user.User, error)
		CreateAddress(
			street,
			streetNumber,
			country,
			city,
			firstName,
			lastName,
			postCode string,
		) user.Address
		CreatePhoneNumber(phoneNumberType common.PhoneNumberType, number string) user.PhoneNumber
		Delete(userID, id string) error
		Get(id string) (*user.User, error)
		GetMany(name *string, sort common.Sort, skip, limit int64) ([]*user.User, error)
		// Update(userID, id string, firstName, lastName *string, isActive *bool, adminUserType *common.AdminUserType, permissions *[]common.RequestPermission, passwordUpdate *PasswordUpdate, locale *common.Locale) (*user.User, error)
	}

	Application struct {
		processLog                        applicationProcessLog.IApplication
		repo                              *repository.Repo
		emailHandler                      emailhandler.IEmailHandler
		landingpageEndpoint               *string
		jwtExpireTimeForPwdResetInMinutes *int
	}
)

func New(mongoDbClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, processLog applicationProcessLog.IApplication, landingpageEndpoint *string, jwtExpireTimeForPwdResetInMinutes *int) IApplication {
	application := Application{
		repo:                              repository.New(mongoDbClient),
		emailHandler:                      emailHandler,
		landingpageEndpoint:               landingpageEndpoint,
		jwtExpireTimeForPwdResetInMinutes: jwtExpireTimeForPwdResetInMinutes,
		processLog:                        processLog,
	}

	return &application
}

func (a *Application) Create(
	email string,
	sex common.Sex,
	birthDate time.Time,
	locale common.Locale,
	firstName string,
	lastName string,
	addresses []user.Address,
	phoneNumbers []user.PhoneNumber,
) (*user.User, error) {
	foundUser, err := a.repo.LoadUserByEmail(email)

	doesRecordExists := err == nil && foundUser != nil
	if doesRecordExists {
		return foundUser, nil
	}

	rndPwd := utils.GenerateRandomString(50)

	newAdminUser := user.New(email, rndPwd, firstName, lastName, birthDate, locale, sex)
	for _, address := range addresses {
		newAdminUser.AddAddress(&address)
	}
	for _, phoneNumber := range phoneNumbers {
		newAdminUser.AddPhoneNumber(&phoneNumber)
	}

	err = a.repo.Insert(newAdminUser)
	if err != nil {
		return nil, fmt.Errorf("%w :%v", ErrCantSave, err)
	}

	// TODO
	// a.processLog.New(userID, email, common.CREATED, common.ADMIN_USER, nil)

	return newAdminUser, nil
}

// func (a *Application) Update(userID, id string, firstName, lastName *string, isActive *bool, adminUserType *common.AdminUserType, permissions *[]common.RequestPermission, passwordUpdate *PasswordUpdate, locale *common.Locale) (*user.User, error) {
// 	foundUser, err := a.repo.LoadAdminUserByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
// 	}

// 	hasPwdUpdate := passwordUpdate != nil
// 	if hasPwdUpdate {
// 		err := foundUser.UpdatePassword(passwordUpdate.OldPassword, passwordUpdate.NewPassword, passwordUpdate.NewComparePassword, false)
// 		if err != nil {
// 			return nil, fmt.Errorf("%w: %v", ErrUpdatingPwdNotPossible, err)
// 		}
// 	}

// 	shouldBeUpdated := permissions != nil
// 	if shouldBeUpdated {
// 		foundUser.Permissions = *permissions
// 	}

// 	err = foundUser.Update(firstName, lastName, isActive, adminUserType, locale)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
// 	}

// 	err = a.repo.Update(userID, foundUser)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v: %s", ErrCantSave, err, id)
// 	}

// 	a.processLog.New(userID, "", common.UPDATED, common.ADMIN_USER, &id)

// 	return foundUser, nil
// }

func (a *Application) Delete(userID, id string) error {
	user, err := a.repo.LoadUserByID(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteUserByID(id)
	if err != nil {
		return fmt.Errorf("%w: %v: %s", ErrCantDelete, err, id)
	}

	a.processLog.New(userID, *user.Email, common.DELETED, common.ADMIN_USER, &id)

	return nil
}

func (a *Application) Get(id string) (*user.User, error) {
	foundUser, err := a.repo.LoadUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundUser, nil
}

func (a *Application) GetMany(name *string, sort common.Sort, skip, limit int64) ([]*user.User, error) {
	foundUsers, err := a.repo.FindMany(name, &sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundUsers, nil
}

func (a *Application) CreateAddress(
	street,
	streetNumber,
	country,
	city,
	firstName,
	lastName,
	postCode string,
) user.Address {
	return user.Address{
		FirstName:    firstName,
		ID:           "1",
		LastName:     lastName,
		StreetName:   street,
		StreetNumber: streetNumber,
		Zip:          postCode,
		City:         city,
		Country:      country,
	}
}

func (a *Application) CreatePhoneNumber(phoneNumberType common.PhoneNumberType, number string) user.PhoneNumber {
	return user.PhoneNumber{
		ID:     "1",
		Type:   &phoneNumberType,
		Number: &number,
	}
}
