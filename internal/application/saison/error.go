package application

import "errors"

var (
	ErrRecordNotExists     = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrSeasonAlreadyExists = errors.New("ErrSeasonAlreadyExists | the season already exists for the year")
	ErrCantUpdate          = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave            = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete          = errors.New("ErrCantDelete | it was not possible to delete the record")
)
