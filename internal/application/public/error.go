package application

import "errors"

var (
	ErrCantSave   = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete = errors.New("ErrCantDelete | it was not possible to delete the record")
	ErrCantUpdate = errors.New("ErrCantUpdate | it was not possible to update the record")

	ErrRecordNotExists = errors.New("the requested record not exists")
)
