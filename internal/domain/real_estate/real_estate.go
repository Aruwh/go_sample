package realEstate

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/pkg/mongodb"
)

type (
	RealEstateLinkedApartments struct {
		ApartmentID   string `json:"id" bson:"_id"`
		ApartmentName string `json:"name" bson:"name"`
	}

	RealEstate struct {
		ID          string                       `json:"id" bson:"_id"`
		Name        *string                      `json:"name" bson:"name"`
		Description shared.Translation           `json:"description" bson:"description"`
		PictureID   *string                      `json:"pictureID" bson:"pictureID"`
		Created     shared.TimeStamp             `json:"created" bson:"created"`
		Edited      shared.TimeStamp             `json:"edited" bson:"edited"`
		Apartments  []RealEstateLinkedApartments `json:"apartments"`
	}
)

func validateInput(name *string, description *shared.Translation) error {
	var error error

	if name == nil {
		error = ErrRealEstateNoName
	}

	if description == nil {
		error = ErrRealEstateNoDescription
	}

	return error
}

func New(creatorID string, name, pictureID *string, description *shared.Translation) (*RealEstate, error) {
	timestamp := shared.NewTimeStamp(&creatorID)

	err := validateInput(name, description)
	if err != nil {
		return nil, err
	}

	realEstate := RealEstate{
		ID:          mongodb.NewID(),
		Name:        name,
		Description: *description,
		PictureID:   pictureID,
		Created:     timestamp,
		Edited:      timestamp,
	}

	return &realEstate, nil
}

func (re *RealEstate) Update(pictureID, name *string, description *shared.Translation) *string {
	var pictureIDToRemove *string

	shouldBeUpdated := pictureID != nil && re.PictureID != pictureID
	if shouldBeUpdated {
		pictureIDToRemove = re.PictureID

		re.PictureID = pictureID
	}

	shouldBeUpdated = name != nil && re.Name != name
	if shouldBeUpdated {
		re.Name = name
	}

	re.Description.Update(description)

	return pictureIDToRemove
}
