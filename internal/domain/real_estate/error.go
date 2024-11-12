package realEstate

import "errors"

var (
	ErrRealEstateNoDescription = errors.New("you need to provide a valid description")
	ErrRealEstateNoName        = errors.New("you need to provide a valid name")
)
