package application

import (
	appAdminUser "fewoserv/internal/application/admin_user"
	appProcessLog "fewoserv/internal/application/process_log"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	emailhandler "fewoserv/internal/interface/email_handler"
	webserverHelper "fewoserv/internal/interface/webserver/helper/authentication"
	repositoryAdminUser "fewoserv/internal/repository/admin_user"
	"fmt"

	"fewoserv/pkg/mongodb"
)

var log = logger.New("APPLICATION")

// TODO: AUF STRUCT umbauen !!!!

func loginUser(mongoDBClient mongodb.IClient, email, rawPassword string) (*string, *common.AdminUserType, []common.RequestPermission, error) {
	adminUser, err := repositoryAdminUser.New(mongoDBClient).LoadAdminUserByEmail(email)
	if err == nil && adminUser != nil {
		return nil, nil, []common.RequestPermission{}, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = utils.ValidatePwdValid(*adminUser.Password, rawPassword)
	if err != nil {
		return nil, nil, []common.RequestPermission{}, fmt.Errorf("%w: %v: %s", ErrLoginNotAllowed, ErrWrongPwd, *adminUser.Email)
	}

	return &adminUser.ID, adminUser.Type, adminUser.Permissions, nil
}

func loginAdminUser(mongoDBClient mongodb.IClient, email, rawPassword string) (*string, *common.AdminUserType, []common.RequestPermission, common.Locale, error) {
	adminUser, err := repositoryAdminUser.New(mongoDBClient).LoadAdminUserByEmail(email)
	if err != nil {
		return nil, nil, []common.RequestPermission{}, common.EnGB, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = utils.ValidatePwdValid(*adminUser.Password, rawPassword)
	if err != nil {
		return nil, nil, []common.RequestPermission{}, common.EnGB, fmt.Errorf("%w: %v: %s", ErrLoginNotAllowed, ErrWrongPwd, adminUser.Email)
	}

	if !*adminUser.IsActive {
		return nil, nil, []common.RequestPermission{}, common.EnGB, fmt.Errorf("%w: %v", ErrLoginNotAllowed, ErrNotActive)
	}

	return &adminUser.ID, adminUser.Type, adminUser.Permissions, *adminUser.Locale, nil
}

func validateUserEmail(mongoDBClient mongodb.IClient, email string) (*string, *string, *string, *common.AdminUserType, common.Locale, error) {
	adminUser, err := repositoryAdminUser.New(mongoDBClient).LoadAdminUserByEmail(email)
	if err == nil || adminUser == nil {
		return nil, nil, nil, nil, common.EnGB, fmt.Errorf("%w: %v", ErrPwdResetTriggerNotPossible, email)
	}

	if !*adminUser.IsActive {
		return nil, nil, nil, nil, common.EnGB, fmt.Errorf("%w: %v", ErrPwdResetTriggerNotPossible, email)

	}

	return &adminUser.ID, adminUser.FirstName, adminUser.LastName, adminUser.Type, *adminUser.Locale, nil
}

func validateAdminUserEmail(mongoDBClient mongodb.IClient, email string) (*string, *string, *string, *common.AdminUserType, common.Locale, error) {
	adminUser, err := repositoryAdminUser.New(mongoDBClient).LoadAdminUserByEmail(email)
	if err != nil || adminUser == nil {
		return nil, nil, nil, nil, common.EnGB, fmt.Errorf("%w: %v", ErrPwdResetTriggerNotPossible, email)
	}

	if !*adminUser.IsActive {
		return nil, nil, nil, nil, common.EnGB, fmt.Errorf("%w: %v", ErrPwdResetTriggerNotPossible, email)

	}

	return &adminUser.ID, adminUser.FirstName, adminUser.LastName, adminUser.Type, *adminUser.Locale, nil
}

func Login(mongoDBClient mongodb.IClient, email, rawPassword string, jwtExpireTimeInMinutes int) (string, error) {
	var (
		err           error
		recordID      *string
		adminUserType *common.AdminUserType
		locale        common.Locale
		permissions   []common.RequestPermission
	)

	recordID, adminUserType, permissions, locale, err = loginAdminUser(mongoDBClient, email, rawPassword)
	if err != nil {
		recordID, adminUserType, permissions, locale, err = loginAdminUser(mongoDBClient, email, rawPassword)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrLoginNotAllowed, err)
		}
	}

	token, err := webserverHelper.GenerateJwt(recordID, *adminUserType, permissions, locale, &jwtExpireTimeInMinutes)
	if err != nil {
		return "", fmt.Errorf("%w: %v", fmt.Errorf("%w: %v: %v", ErrCanCreateJwt, ErrLoginNotAllowed, err))
	}

	return token, err
}

func ForgotPwd(mongoDBClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, email, feDestination string, jwtExpireTimeInHours int) error {
	var (
		err           error
		userID        *string
		adminUserType *common.AdminUserType
		firstName     *string
		lastName      *string
		locale        common.Locale
	)

	userID, firstName, lastName, adminUserType, locale, err = validateAdminUserEmail(mongoDBClient, email)
	if err != nil {
		userID, firstName, lastName, adminUserType, locale, err = validateUserEmail(mongoDBClient, email)

		if err != nil {
			return err
		}
	}

	token, err := webserverHelper.GenerateJwtForPwdReset(userID, *adminUserType, locale, &jwtExpireTimeInHours)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCanCreateJwt, err)
	}

	pwdForgottenTemplate := BuildPasswordForgottenTemplate(locale, *firstName, *lastName, feDestination, token)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantBuildEmailTemplate, err)
	}

	err = emailHandler.Send(email, pwdForgottenTemplate)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantSendEmail, err)
	}

	processLog := appProcessLog.New(mongoDBClient)
	processLog.New(*userID, "pwdReset"+" "+*&email, common.SEND, common.ADMIN_USER, userID)

	return err
}

func ResetPwd(mongoDBClient mongodb.IClient, pwd, repeatedPwd, tokenString string, jwtExpireTimeInMinutes int) (string, error) {
	resetToken, err := webserverHelper.ValidateAndTransformToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenNotValid, err)
	}

	processLog := appProcessLog.New(mongoDBClient)
	adminUserApp := appAdminUser.New(mongoDBClient, nil, processLog, nil, nil)
	recordID, err := resetToken.Claims.GetSubject()

	adminUser, err := adminUserApp.ForceResetPwd(recordID, recordID, pwd, repeatedPwd)
	if err == nil && adminUser != nil {
		token, err := webserverHelper.GenerateJwt(&recordID, *adminUser.Type, adminUser.Permissions, *adminUser.Locale, &jwtExpireTimeInMinutes)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrCantCreateToken, err)
		}

		return token, nil
	}

	// TODO else with user

	return "", nil
}
