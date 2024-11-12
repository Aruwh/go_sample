package application

import "errors"

var (
	ErrRecordNotExists        = errors.New("ErrRecordNotExists | the requested record not exists")
	ErrNotAllowedDeletion     = errors.New("ErrNotAllowedDeletion | it is not allowed to delete the super admin")
	ErrCantCreateToken        = errors.New("ErrCantCreateToken | it ws not possbile to create a token")
	ErrCantSendInvitation     = errors.New("ErrCantSendInvitation | it ws not possbile to send the invitation email")
	ErrPwdNotEqual            = errors.New("ErrPwdNotEqual | the provided passwords aren't equal")
	ErrCantInitAdminUser      = errors.New("ErrCantInitAdminUser | it was not possbile to initiate a adminUser")
	ErrUpdatingPwdNotPossible = errors.New("ErrUpdatingPwdNotPossible | something went wrong during the pwd updade")
	ErrCantUpdatePwd          = errors.New("ErrCantUpdatePwd | it was not possible to update the password")

	ErrCantUpdate = errors.New("ErrCantUpdate | it was not possible to update the record")
	ErrCantSave   = errors.New("ErrCantSave | it ws not possbile to write the record to the DB")
	ErrCantDelete = errors.New("ErrCantDelete | it was not possible to delete the record")

	// not for external usage (FE)
	ErrEmailAlreadyInSystem = errors.New("ErrEmailAlreadyInSystem | the provided email is already in the system")
)
