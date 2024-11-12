package user

import "errors"

var (
	ErrUserPasswordMissMatch       = errors.New("the provided passwords are not equal")
	ErrUserNewPasswordsAreNotEqual = errors.New("the provided new passwords aren't equal")
)
