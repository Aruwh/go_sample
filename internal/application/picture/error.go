package application

import "errors"

var (
	ErrRecordNotExists       = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrCantCreateStorageDir  = errors.New("ErrCantCreateStorageDir | it was not possible to create the dir")
	ErrCantStoreFile         = errors.New("ErrCantStoreFile | it was not possible to store the file on the path")
	ErrCantDeleteFile        = errors.New("ErrCantDeleteFile | it was not possible to delete the file from path")
	ErrCantTransformToBase64 = errors.New("ErrCantTransformToBase64 | it was not possible transform the picture to base64")
	ErrCantUpdate = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave   = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete = errors.New("ErrCantDelete | it was not possible to delete the record")
)
