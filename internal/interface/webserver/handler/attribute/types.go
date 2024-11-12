package attribute

import (
	"fewoserv/internal/interface/webserver/shared"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Name *shared.Translation `json:"name"`
		Svg  *string             `json:"svg"`
	}

	UpdateRequest struct {
		Name *shared.Translation `json:"name"`
		Svg  *string             `json:"svg"`
	}

	GetManyFilter struct {
		Name *string `json:"name"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
)
