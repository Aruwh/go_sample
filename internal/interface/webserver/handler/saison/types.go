package saison

import (
	"fewoserv/internal/interface/webserver/shared"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Year   int                  `json:"year"`
		Enries []shared.SaisonEntry `json:"entries"`
	}

	UpdateRequest struct {
		Year   *int                  `json:"year"`
		Enries *[]shared.SaisonEntry `json:"entries"`
	}

	GetManyFilter struct {
		Year *int `json:"year"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

)
