package application

import "errors"

var (
	ErrRecordNotExists      = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrCantIncBookingNumber = errors.New("ErrCantIncBookingNumber | it was not possible to increment the booking number")
	ErrCantUpdate           = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave             = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete           = errors.New("ErrCantDelete | it was not possible to delete the record")
)
