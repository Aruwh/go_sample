package shared

import "errors"

var (
	ErrUnsupportedAdminType = errors.New("the provided type is not supported")
)
