package application

import "errors"

var (
	ErrRecordNotExists = errors.New("the requested record not exists")
	ErrCantUnmarshal   = errors.New("it was not possible to unmarshal the json")

	ErrCantUpdate = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave   = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete = errors.New("ErrCantDelete | it was not possible to delete the record")
)
