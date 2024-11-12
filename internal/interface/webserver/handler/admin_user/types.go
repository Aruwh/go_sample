package adminuser

import (
	"fewoserv/internal/infrastructure/common"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Email     string               `json:"email"`
		IsActive  bool                 `json:"isActive"`
		FirstName string               `json:"firstName" validate:"min=2"`
		LastName  string               `json:"lastName" validate:"min=2"`
		Type      common.AdminUserType `json:"type"`
		Locale    common.Locale        `json:"locale"`
	}

	PasswordUpdate struct {
		OldPassword        *string `json:"oldPassword"`
		NewPassword        *string `json:"newPassword" validate:"min=6"`
		NewComparePassword *string `json:"newComparePassword" validate:"min=6"`
	}

	UpdateRequest struct {
		Type           *common.AdminUserType       `json:"type"`
		Locale         *common.Locale              `json:"locale"`
		FirstName      *string                     `json:"firstName"`
		LastName       *string                     `json:"lastName"`
		IsActive       *bool                       `json:"isActive"`
		PasswordUpdate *PasswordUpdate             `json:"passwordUpdate"`
		Permissions    *[]common.RequestPermission `json:"permissions"`
	}

	GetManyFilter struct {
		Type *common.AdminUserType `json:"type"`
		Name *string               `json:"name"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

)
