package booking

import "errors"

var (
	ErrDateUpdateNotAllowed   = errors.New("you need to provide the from and the to date")
	ErrStatusSwitchNotAllowed = errors.New("it is not allowed to switch the status")
)
