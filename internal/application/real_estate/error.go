package application

import "errors"

var (
	ErrCantInitRealestate = errors.New("ErrCantInitRealestate | it was not possible to init a realEstate")
	ErrRecordNotExists     = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrHasLinkedApartments = errors.New("ErrHasLinkedApartments | the real estate can't be deleted, because of linked apartment")
	ErrCantUpdate = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave   = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete = errors.New("ErrCantDelete | it was not possible to delete the record")
)
