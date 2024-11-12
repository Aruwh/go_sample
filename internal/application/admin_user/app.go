package application

import (
	applicationProcessLog "fewoserv/internal/application/process_log"
	adminuser "fewoserv/internal/domain/admin_user"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	emailhandler "fewoserv/internal/interface/email_handler"
	webserverHelper "fewoserv/internal/interface/webserver/helper/authentication"
	repository "fewoserv/internal/repository/admin_user"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

var (
	log                            = logger.New("APPLICATION")
	defaultAdminRequestPermissions = []common.RequestPermission{
		common.ADMIN_USER_VIEW,
		common.ADMIN_USER_EDIT,
		common.ADMIN_USER_DELETE,
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
	defaultApartmentOwnerRequestPermissions = []common.RequestPermission{
		common.APARTMENT_VIEW,
		common.BOOKING_VIEW,
		common.BOOKING_EDIT,
		common.SETTINGS_VIEW,
		common.PICTURE_VIEW,
	}
)

type (
	PasswordUpdate struct {
		OldPassword        string
		NewPassword        string
		NewComparePassword string
	}

	IApplication interface {
		CreateSuperAdmin() (*adminuser.AdminUser, error)
		CreateAdmin(userID, email, rawPassword, compareRawPassword, firstName, lastName string, locale common.Locale) (*adminuser.AdminUser, error)
		CreateApartmentOwner(userID, email, rawPassword, compareRawPassword, firstName, lastName string, locale common.Locale) (*adminuser.AdminUser, error)
		Create(userID, email, rawPassword, compareRawPassword, firstName, lastName string, permissions []common.RequestPermission, adminType common.AdminUserType, locale common.Locale) (*adminuser.AdminUser, error)
		Delete(userID, id string) error
		Get(id string) (*adminuser.AdminUser, error)
		GetMany(adminUserType *common.AdminUserType, name *string, sort common.Sort, skip, limit int64) ([]*adminuser.AdminUser, error)
		Update(userID, id string, firstName, lastName *string, isActive *bool, adminUserType *common.AdminUserType, permissions *[]common.RequestPermission, passwordUpdate *PasswordUpdate, locale *common.Locale) (*adminuser.AdminUser, error)
		SendInvitation(userID, id string) error
		ForceResetPwd(userID, id, pwd, repeatedPwd string) (*adminuser.AdminUser, error)
	}

	Application struct {
		processLog                        applicationProcessLog.IApplication
		repo                              *repository.Repo
		emailHandler                      emailhandler.IEmailHandler
		feEndpoint                        *string
		jwtExpireTimeForPwdResetInMinutes *int
	}
)

func New(mongoDbClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, processLog applicationProcessLog.IApplication, feEndpoint *string, jwtExpireTimeForPwdResetInMinutes *int) IApplication {
	application := Application{
		repo:                              repository.New(mongoDbClient),
		emailHandler:                      emailHandler,
		feEndpoint:                        feEndpoint,
		jwtExpireTimeForPwdResetInMinutes: jwtExpireTimeForPwdResetInMinutes,
		processLog:                        processLog,
	}

	return &application
}

func (a *Application) SendInvitation(userID, id string) error {
	foundAdminUser, err := a.repo.LoadAdminUserByID(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	token, err := webserverHelper.GenerateJwtForPwdReset(&foundAdminUser.ID, *foundAdminUser.Type, *foundAdminUser.Locale, a.jwtExpireTimeForPwdResetInMinutes)
	if err != nil {
		return fmt.Errorf("%w :%v", ErrCantCreateToken, err)
	}

	invitationTemplate := BuildEmailInvitationTemplate(*foundAdminUser.Locale, *foundAdminUser.FirstName, *foundAdminUser.LastName, *a.feEndpoint, token)

	err = a.emailHandler.Send(*foundAdminUser.Email, invitationTemplate)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantSendInvitation, err)
	}

	// set information regarding invitation
	timeStamp := shared.NewTimeStamp(&userID)
	foundAdminUser.Invited.At = &timeStamp
	foundAdminUser.Invited.AcceptedAt = nil

	err = a.repo.Update(userID, foundAdminUser)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	a.processLog.New(userID, "invitation"+" "+*foundAdminUser.Email, common.SEND, common.ADMIN_USER, &id)

	return nil
}

func (a *Application) CreateSuperAdmin() (*adminuser.AdminUser, error) {
	fmt.Println("start creating super admin user:")

	newSuperAdmin, plainPassword, err := adminuser.NewSuperUser(nil)
	if err != nil {
		fmt.Println("CreateSuperAdmin :: Error:", err)
		return nil, err
	}

	newSuperAdmin.AddPermissions(defaultAdminRequestPermissions...)

	existingAdminUser, _ := a.repo.LoadAdminUserByEmail(*newSuperAdmin.Email)
	if existingAdminUser != nil {
		fmt.Println("CreateSuperAdmin :: Admin user already exists !")
		return existingAdminUser, nil
	}

	err = a.repo.Insert(newSuperAdmin)
	if err == nil {
		fmt.Println("----------------------------------------------------------------------------------------")
		fmt.Println("Super Admin created  ")
		fmt.Println("email: ", *newSuperAdmin.Email)
		fmt.Println("password: ", plainPassword)
		fmt.Println("")
		fmt.Println("save the password in a safe place, otherwise it will be gone forever !")
		fmt.Println("----------------------------------------------------------------------------------------")

		return newSuperAdmin, nil
	}

	return nil, err
}

func (a *Application) CreateAdmin(userID, email, rawPassword, compareRawPassword, firstName, lastName string, locale common.Locale) (*adminuser.AdminUser, error) {
	newAdmin, err := a.Create(userID, email, rawPassword, compareRawPassword, firstName, lastName, defaultAdminRequestPermissions, common.ADMINISTRATOR, locale)
	if err != nil {
		return nil, err
	}

	return newAdmin, err
}

func (a *Application) CreateApartmentOwner(userID, email, rawPassword, compareRawPassword, firstName, lastName string, locale common.Locale) (*adminuser.AdminUser, error) {
	newAdmin, err := a.Create(userID, email, rawPassword, compareRawPassword, firstName, lastName, defaultApartmentOwnerRequestPermissions, common.APARTMENT_OWNER, locale)
	if err != nil {
		return nil, err
	}

	return newAdmin, err
}

func (a *Application) Create(userID, email, rawPassword, compareRawPassword, firstName, lastName string, permissions []common.RequestPermission, adminType common.AdminUserType, locale common.Locale) (*adminuser.AdminUser, error) {
	isPasswordEqual := rawPassword == compareRawPassword
	if !isPasswordEqual {
		return nil, ErrPwdNotEqual
	}

	newAdminUser, err := adminuser.NewUser(userID, email, rawPassword, compareRawPassword, firstName, lastName, adminType, locale)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantInitAdminUser, err)
	}

	newAdminUser.AddPermissions(permissions...)

	isExisting := a.repo.ValidateAdminUserExists(newAdminUser)
	if !isExisting {
		a.repo.Insert(newAdminUser)

		a.processLog.New(userID, email, common.CREATED, common.ADMIN_USER, nil)

		return newAdminUser, nil
	}

	return nil, fmt.Errorf("%w: %v : %s", ErrRecordNotExists, ErrEmailAlreadyInSystem, email)
}

func (a *Application) Update(userID, id string, firstName, lastName *string, isActive *bool, adminUserType *common.AdminUserType, permissions *[]common.RequestPermission, passwordUpdate *PasswordUpdate, locale *common.Locale) (*adminuser.AdminUser, error) {
	foundAdminUser, err := a.repo.LoadAdminUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	hasPwdUpdate := passwordUpdate != nil
	if hasPwdUpdate {
		err := foundAdminUser.UpdatePassword(passwordUpdate.OldPassword, passwordUpdate.NewPassword, passwordUpdate.NewComparePassword, false)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrUpdatingPwdNotPossible, err)
		}
	}

	shouldBeUpdated := permissions != nil
	if shouldBeUpdated {
		foundAdminUser.Permissions = *permissions
	}

	err = foundAdminUser.Update(firstName, lastName, isActive, adminUserType, locale)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	err = a.repo.Update(userID, foundAdminUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %v: %s", ErrCantSave, err, id)
	}

	a.processLog.New(userID, "", common.UPDATED, common.ADMIN_USER, &id)

	return foundAdminUser, nil
}

func (a *Application) Delete(userID, id string) error {
	adminUser, err := a.repo.LoadAdminUserByID(id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	deletionIsNotAllowed := *adminUser.Type == common.ADMINISTRATOR
	if deletionIsNotAllowed {
		return ErrNotAllowedDeletion
	}

	err = a.repo.DeleteAdminUserByID(id)
	if err != nil {
		return fmt.Errorf("%w: %v: %s", ErrCantDelete, err, id)
	}

	a.processLog.New(userID, *adminUser.Email, common.DELETED, common.ADMIN_USER, &id)

	return nil
}

func (a *Application) Get(id string) (*adminuser.AdminUser, error) {
	foundAdminUser, err := a.repo.LoadAdminUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundAdminUser, nil
}

func (a *Application) GetMany(adminUserType *common.AdminUserType, name *string, sort common.Sort, skip, limit int64) ([]*adminuser.AdminUser, error) {
	foundAdminUsers, err := a.repo.FindMany(adminUserType, name, &sort, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	return foundAdminUsers, nil
}

func (a *Application) ForceResetPwd(userID, id, pwd, repeatedPwd string) (*adminuser.AdminUser, error) {
	foundAdminUser, err := a.repo.LoadAdminUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = foundAdminUser.UpdatePassword("", pwd, repeatedPwd, true)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdatePwd, err)
	}

	// the first time when the pwd was reset was a trigger from the invitation process
	isFirstPWDReset := foundAdminUser.Invited.AcceptedAt == nil
	if isFirstPWDReset {
		acceptedAt := time.Now()
		foundAdminUser.Invited.AcceptedAt = &acceptedAt

		a.processLog.New(userID, "invitation"+" "+*foundAdminUser.Email, common.ACCEPTED, common.ADMIN_USER, &id)
	}

	err = a.repo.Update(userID, foundAdminUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	a.processLog.New(userID, "", common.UPDATED, common.ADMIN_USER, &id)

	return foundAdminUser, nil
}
