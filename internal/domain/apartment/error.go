package apartment

import "errors"

var (
	ErrRealEstateNoDescription = errors.New("you need to provide a valid description")
	ErrUnknownSaisonType       = errors.New("the provided saison type is unknown")
)
