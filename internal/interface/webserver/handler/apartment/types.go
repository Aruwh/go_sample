package apartment

import (
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/interface/webserver/shared"
)

type (

	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	CreateRequest struct {
		Name                  string                 `json:"name"`
		OwnerID               string                 `json:"ownerID" validate:"mongoDbID"`
		RealEstateID          string                 `json:"realEstateID" validate:"mongoDbID"`
		Description           *shared.Translation    `json:"description"`
		PictureIDs            *[]string              `json:"pictureIDs"`
		SaisonPrice           *apartment.SaisonPrice `json:"saisonPrice"`
		AttributeIDs          *[]string              `json:"attributeIDs"`
		TopAttributeIDs       *[]string              `json:"topAttributeIDs"`
		RoomSize              int                    `json:"roomSize"`
		SleepingPlaces        int                    `json:"sleepingPlaces"`
		Bathrooms             int                    `json:"bathRooms"`
		AllowedNumberOfPeople int                    `json:"allowedNumberOfPeople"`
		AllowedNumberOfPets   int                    `json:"allowedNumberOfPets"`
	}

	UpdateRequest struct {
		Name                  *string                `json:"name"`
		IsActive              *bool                  `json:"isActive"`
		OwnerID               *string                `json:"ownerID"`
		RealEstateID          *string                `json:"realEstateID"`
		Description           *shared.Translation    `json:"description"`
		PictureIDs            *[]string              `json:"pictureIDs"`
		SaisonPrice           *apartment.SaisonPrice `json:"saisonPrice"`
		AttributeIDs          *[]string              `json:"attributeIDs"`
		TopAttributeIDs       *[]string              `json:"topAttributeIDs"`
		RoomSize              *int                   `json:"roomSize"`
		Bathrooms             *int                    `json:"bathRooms"`
		SleepingPlaces        *int                   `json:"sleepingPlaces"`
		AllowedNumberOfPeople *int                   `json:"allowedNumberOfPeople"`
		AllowedNumberOfPets   *int                    `json:"allowedNumberOfPets"`
	}

	GetManyFilter struct {
		RealEstateID *string `json:"realEstateID"`
		Name         *string `json:"name"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

	GetManyPublicResponse struct {
		Name                string                        `json:"name"`
		Description         interface{}                   `json:"description"`
		PictureIDs          []string                      `json:"pictureIDs"`
		AttributeCollection apartment.AttributeCollection `json:"attributeCollection"`
		SaisonPrice         apartment.SaisonPrice         `json:"saisonPrice"`
	}
)

func NewGetManyPublicResponse(apartments []*apartment.Apartment) []GetManyPublicResponse {
	reducedApartments := []GetManyPublicResponse{}

	for _, apartment := range apartments {
		reducedApartment := GetManyPublicResponse{
			Name:                *apartment.Name,
			Description:         *apartment.Description,
			PictureIDs:          apartment.PictureIDs,
			AttributeCollection: *apartment.AttributeCollection,
			SaisonPrice:         *apartment.SaisonPrice,
		}

		reducedApartments = append(reducedApartments, reducedApartment)
	}

	return reducedApartments
}
