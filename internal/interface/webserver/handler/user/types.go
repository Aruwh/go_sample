package user

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/shared"
	"time"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Email        string               `json:"email" bson:"email"`
		Password     string               `json:"password" bson:"password"`
		Sex          common.Sex           `json:"sex" bson:"sex"`
		BirthDate    time.Time            `json:"birthDate" bson:"birthDate"`
		Locale       common.Locale        `json:"locale" bson:"locale"`
		FirstName    string               `json:"firstName" bson:"firstName"`
		LastName     string               `json:"lastName" bson:"lastName"`
		Addresses    []shared.Address     `json:"addresses" bson:"addresses"`
		PhoneNumbers []shared.PhoneNumber `json:"phoneNumbers" bson:"phoneNumbers"`
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
		Name *string `json:"name"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

)
