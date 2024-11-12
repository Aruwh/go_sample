package picture

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/interface/webserver/shared"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	UploadRequest struct {
		Description *shared.Translation   `json:"description"`
		Variant     common.PictureVariant `json:"variant"`
	}

	UpdateRequest struct {
		Description *shared.Translation   `json:"description"`
		Variant     common.PictureVariant `json:"variant"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
)
