package application

import "errors"

var (
	ErrIncBookingNumber         = errors.New("ErrIncBookingNumber | it was not possible to increase the booking number")
	ErrRecordNotExists          = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrPlacingBookingNotAllowed = errors.New("ErrPlacingBookingNotAllowed | it was not possible to place the booking, because it overlaps with a other one")
	ErrCantUpdate               = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave                 = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete               = errors.New("ErrCantDelete | it was not possible to delete the record")
)
