package application

import "errors"

var (
	ErrLoginNotAllowed        = errors.New("ErrLoginNotAllowed | the user is not allowed to login")
	ErrCanCreateJwt           = errors.New("ErrCanCreateJwt | it was not possible to create a jwt")
	ErrCantBuildEmailTemplate = errors.New("ErrCantBuildEmailTemplate | it was not possible to create the email template")
	ErrCantSendEmail          = errors.New("ErrCantSendEmail | it was not possible to send the email")
	ErrTokenNotValid          = errors.New("ErrTokenNotValid | the provided token is not valid")
	ErrCantCreateToken        = errors.New("ErrCantCreateToken | it was not possible to create the token")
	ErrRecordNotExists            = errors.New("the requested record not exists")
	ErrWrongPwd                   = errors.New("the given pwd is wrong")
	ErrNotActive                  = errors.New("the user is not active")
	ErrPwdResetTriggerNotPossible = errors.New("it is not possible to trigger the pwd reset")
)
