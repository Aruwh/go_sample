package apartment

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/pkg/mongodb"
)

var logApartment = logger.New("APARTMENT")

type (
	Apartment struct {
		ID                  string               `json:"id" bson:"_id"`
		IsActive            bool                 `json:"isActive" bson:"isActive"`
		OwnerID             *string              `json:"ownerID" bson:"ownerID"`
		RealEstateID        *string              `json:"realEstateID" bson:"realEstateID"`
		Name                *string              `json:"name" bson:"name"`
		Description         *shared.Translation  `json:"description" bson:"description"`
		PictureIDs          []string             `json:"pictureIDs" bson:"pictureIDs"`
		AttributeCollection *AttributeCollection `json:"attributeCollection" bson:"attributeCollection"`
		SaisonPrice         *SaisonPrice         `json:"saisonPrice" bson:"saisonPrice"`
		Created             *shared.TimeStamp    `json:"created" bson:"created"`
		Edited              *shared.TimeStamp    `json:"edited" bson:"edited"`
	}

	MinimalApartment struct {
		ID             string              `json:"id" bson:"_id"`
		IsActive       bool                `json:"isActive" bson:"isActive"`
		RealEstateID   *string             `json:"realEstateID" bson:"realEstateID"`
		Name           *string             `json:"name" bson:"name"`
		RealEstateName *string             `json:"realEstateName" bson:"realEstateName"`
		Description    *shared.Translation `json:"description" bson:"description"`
		PictureIDs     []string            `json:"pictureIDs" bson:"pictureIDs"`
	}

	ReadOnlyApartment struct {
		Name                *string                      `json:"name" bson:"name"`
		Description         *shared.Translation          `json:"description" bson:"description"`
		PictureIDs          []string                     `json:"pictureIDs" bson:"pictureIDs"`
		AttributeCollection *AttributeCollectionReadOnly `json:"attributeCollection" bson:"attributeCollection"`
		SaisonPrice         *SaisonPrice                 `json:"saisonPrice" bson:"saisonPrice"`
	}
)

func New(createrID, OwnerID, realEstateID, name string) *Apartment {
	apartmentID := mongodb.NewID()
	timeStamp := shared.NewTimeStamp(&createrID)

	var apartment = Apartment{
		ID:                  apartmentID,
		IsActive:            false,
		OwnerID:             &OwnerID,
		RealEstateID:        &realEstateID,
		Name:                &name,
		Description:         shared.NewTranslation(""),
		PictureIDs:          []string{},
		AttributeCollection: NewAttributeCollection(),
		SaisonPrice:         NewSaisonPrice(nil, nil, nil, nil, nil, nil),
		Created:             &timeStamp,
		Edited:              &timeStamp,
	}

	return &apartment
}

func (a *Apartment) UpdateRealEstateID(realEstateID string) {
	shouldBeUpdated := realEstateID != *a.RealEstateID
	if shouldBeUpdated {
		a.RealEstateID = &realEstateID
	}
}

func (a *Apartment) UpdateName(name *string) {
	shouldBeUpdated := name != nil && name != a.Name
	if shouldBeUpdated {
		a.Name = name
	}
}

func (a *Apartment) UpdateDescription(description *shared.Translation) {
	a.Description.Update(description)
}

func (a *Apartment) UpdatePictureIDs(pictureIDsToUpdate *[]string) []string {
	canIDoSomething := pictureIDsToUpdate != nil
	if !canIDoSomething {
		return []string{}
	}

	addedIDs := utils.Intersection[string](*pictureIDsToUpdate, a.PictureIDs)
	idsToRemove := utils.Intersection[string](a.PictureIDs, *pictureIDsToUpdate)

	updatedPictureIDs := []string{}

	if len(idsToRemove) == 0 {
		updatedPictureIDs = a.PictureIDs
	}

	for _, idToRemove := range idsToRemove {
		for _, pictureID := range a.PictureIDs {
			shouldNotBeRemoved := pictureID != idToRemove
			if shouldNotBeRemoved {
				updatedPictureIDs = append(updatedPictureIDs, pictureID)
			}
		}
	}

	if len(addedIDs) != 0 {
		updatedPictureIDs = append(updatedPictureIDs, addedIDs...)
	}

	if len(updatedPictureIDs) != 0 {
		a.PictureIDs = updatedPictureIDs
	}

	a.PictureIDs = *pictureIDsToUpdate

	return idsToRemove
}

func (a *Apartment) UpdateSaisonPrice(saisonPrice *SaisonPrice) {
	if saisonPrice == nil {
		return
	}

	a.SaisonPrice = saisonPrice
}

func (a *Apartment) UpdateAttributeIDs(attributeIDs *[]string) {
	canIDoSomething := attributeIDs != nil
	if !canIDoSomething {
		return
	}

	a.AttributeCollection.UpdateAttributeIDs(*attributeIDs...)
}

func (a *Apartment) UpdateTopAttributeIDs(topAttributeIDs *[]string) {
	canIDoSomething := topAttributeIDs != nil
	if !canIDoSomething {
		return
	}

	a.AttributeCollection.UpdateTopAttributeIDs(*topAttributeIDs...)
}
