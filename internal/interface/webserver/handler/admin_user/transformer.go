package adminuser

import (
	application "fewoserv/internal/application/admin_user"
)

func TransformRequestUpdatePassword(passwordUpdate *PasswordUpdate) *application.PasswordUpdate {
	if passwordUpdate == nil {
		return nil
	}

	transformedPasswordUpdate := application.PasswordUpdate{
		OldPassword:        *passwordUpdate.OldPassword,
		NewPassword:        *passwordUpdate.NewPassword,
		NewComparePassword: *passwordUpdate.NewComparePassword,
	}

	return &transformedPasswordUpdate
}
