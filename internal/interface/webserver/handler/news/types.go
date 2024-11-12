package news

import (
	"fewoserv/internal/interface/webserver/shared"
	"time"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Title     shared.Translation `json:"title"`
		Content   shared.Translation `json:"content"`
		PublishAt time.Time          `json:"publishAt"`
	}

	UpdateRequest struct {
		Title     shared.Translation `json:"title"`
		Content   shared.Translation `json:"content"`
		PublishAt time.Time          `json:"publishAt"`
		Active    *bool              `json:"active"`
	}

	GetManyFilter struct {
		Title *string `json:"title"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
)
