package adminuser

import "errors"

var (
	ErrInvalidEmail                  = errors.New("ErrInvalidEmail | the passed email is not valid")
	ErrPasswordsNotEqual             = errors.New("ErrPasswordsNotEqual | the passwords aren't equal")
	ErrInvalidPasswords              = errors.New("ErrInvalidPasswords | the passwords is not valid")
	ErrFirstNameToShort              = errors.New("ErrFirstNameToShort | the passed firstName is to short")
	ErrLastNameToShort               = errors.New("ErrLastNameToShort | the passed lastName is to short")
	ErrAdminUserTypeChangeNotAllowed = errors.New("ErrAdminUserTypeChangeNotAllowed | the change is not allowed")
	ErrUserPasswordMissMatch         = errors.New("ErrUserPasswordMissMatch | the provided passwords are not equal")
	ErrUserNewPasswordsAreNotEqual   = errors.New("ErrUserNewPasswordsAreNotEqual | the provided new passwords aren't equal")
)
