package realestate

import "fewoserv/internal/interface/webserver/shared"

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Name        string             `json:"name"`
		Description shared.Translation `json:"description"`
		PictureID   *string            `json:"pictureID"`
	}

	UpdateRequest struct {
		Name        *string             `json:"name"`
		Description *shared.Translation `json:"description"`
		PictureID   *string             `json:"pictureID"`
	}

	GetManyFilter struct {
		Name *string `json:"name"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
)
