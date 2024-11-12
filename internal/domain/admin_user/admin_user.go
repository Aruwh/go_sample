package adminuser

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/pkg/mongodb"
	"time"
)

var log = logger.New("ADMIN_USER")

type (
	Invitation struct {
		At         *shared.TimeStamp `json:"at" bson:"at"`
		AcceptedAt *time.Time        `json:"acceptedAt" bson:"acceptedAt"`
	}

	AdminUser struct {
		ID          string                     `json:"id" bson:"_id"`
		Email       *string                    `json:"email" bson:"email"`
		Password    *string                    `json:"-" bson:"password"`
		IsActive    *bool                      `json:"isActive" bson:"isActive"`
		FirstName   *string                    `json:"firstName" bson:"firstName"`
		LastName    *string                    `json:"lastName" bson:"lastName"`
		Type        *common.AdminUserType      `json:"type" bson:"type"`
		Permissions []common.RequestPermission `json:"permissions" bson:"permissions"`
		Locale      *common.Locale             `json:"locale" bson:"locale"`
		Invited     Invitation                 `json:"invited" bson:"invited"`
		Created     shared.TimeStamp           `json:"created" bson:"created"`
		Edited      shared.TimeStamp           `json:"edited" bson:"edited"`
	}
)

func NewUser(creatorID, email, rawPassword, compareRawPassword, firstName, lastName string, adminUserType common.AdminUserType, locale common.Locale) (*AdminUser, error) {
	var (
		usedType  = adminUserType
		userID    = mongodb.NewID()
		timeStamp = shared.NewTimeStamp(&creatorID)
	)

	decryptedPassword, err := utils.DecryptPwd(rawPassword)
	if err != nil {
		return nil, err
	}

	userIsActive := true

	var adminUser = AdminUser{
		ID:          userID,
		IsActive:    &userIsActive,
		Email:       &email,
		Password:    decryptedPassword,
		FirstName:   &firstName,
		LastName:    &lastName,
		Type:        &usedType,
		Permissions: []common.RequestPermission{},
		Locale:      &locale,
		Invited:     Invitation{},
		Created:     timeStamp,
		Edited:      timeStamp,
	}

	return &adminUser, nil
}

func NewSuperUser(creatorID *string) (*AdminUser, string, error) {
	email := common.SuperAdminEmail
	password := utils.GenerateRandomString(24)
	firstName := "Super"
	lastName := "AdminUser"

	usedCreatorID := mongodb.NewID()
	if creatorID != nil {
		usedCreatorID = *creatorID
	}
	superAdminUser, err := NewUser(usedCreatorID, email, password, password, firstName, lastName, common.SUPER_ADMINISTRATOR, common.EnGB)

	return superAdminUser, password, err
}

func (au *AdminUser) AddPermissions(permissions ...common.RequestPermission) {
	au.Permissions = append(au.Permissions, permissions...)
}

func (au *AdminUser) RemovePermissions(permissionsToRemove ...common.RequestPermission) {
	var reducedPermissions = []common.RequestPermission{}

	for _, permission := range au.Permissions {
		for _, permissionToRemove := range permissionsToRemove {
			if permissionToRemove != permission {
				reducedPermissions = append(reducedPermissions, permission)
			}
		}
	}

	au.Permissions = reducedPermissions
}

func (u *AdminUser) Update(firstName, lastName *string, isActive *bool, adminUserType *common.AdminUserType, locale *common.Locale) error {
	shouldBeUpdated := firstName != nil && *u.FirstName != *firstName
	if shouldBeUpdated {
		u.FirstName = firstName
	}

	shouldBeUpdated = lastName != nil && *u.LastName != *lastName
	if shouldBeUpdated {
		u.LastName = lastName
	}

	shouldBeUpdated = adminUserType != nil && *u.Type != *adminUserType
	if shouldBeUpdated {
		changeIsNotAllowed := *adminUserType == common.SUPER_ADMINISTRATOR
		if changeIsNotAllowed {
			return ErrAdminUserTypeChangeNotAllowed
		}
		u.Type = adminUserType
	}

	shouldBeUpdated = isActive != nil && *u.IsActive != *isActive
	if shouldBeUpdated {
		u.IsActive = isActive
	}

	shouldBeUpdated = locale != nil && *u.Locale != *locale
	if shouldBeUpdated {
		u.Locale = locale
	}

	return nil
}

func (u *AdminUser) UpdatePassword(oldPassword, newPassword, repeatedNewPassword string, ignoreCompareValidation bool) error {
	if ignoreCompareValidation == false {
		err := utils.ValidatePwdValid(*u.Password, oldPassword)
		if err != nil {
			return ErrUserPasswordMissMatch
		}
	}

	areNewPasswordsEqual := newPassword == repeatedNewPassword
	if !areNewPasswordsEqual {
		return ErrUserNewPasswordsAreNotEqual
	}

	decryptedNewPassword, err := utils.DecryptPwd(newPassword)
	if err != nil {
		return err
	}

	u.Password = decryptedNewPassword

	return nil
}
