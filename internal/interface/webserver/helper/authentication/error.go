package helper

import "errors"

var (
	ErrNotAuthorised  = errors.New("you are not authorized to perform this operation")
	ErrMalformedToken = errors.New("the provided token is malformed")
	ErrNoPermission   = errors.New("you don't have the permission to perform this operation")
)
