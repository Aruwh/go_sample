package apartment

import "fewoserv/internal/domain/shared"

type (
	AttributeCollection struct {
		RoomSize              int      `json:"roomSize" bson:"roomSize"`
		SleepingPlaces        int      `json:"sleepingPlaces" bson:"sleepingPlaces"`
		Bathrooms             int      `json:"bathRooms" bson:"bathRooms"`
		AllowedNumberOfPeople int      `json:"allowedNumberOfPeople" bson:"allowedNumberOfPeople"`
		AllowedNumberOfPets   int      `json:"allowedNumberOfPets" bson:"allowedNumberOfPets"`
		TopAttributeIDs       []string `json:"topAttributeIDs" bson:"topAttributeIDs"`
		AttributeIDs          []string `json:"attributeIDs" bson:"attributeIDs"`
	}

	AttributeReadOnly struct {
		Name shared.Translation `json:"name" bson:"name"`
		Svg  string             `json:"svg" bson:"svg"`
	}

	AttributeCollectionReadOnly struct {
		RoomSize              int                 `json:"roomSize" bson:"roomSize"`
		SleepingPlaces        int                 `json:"sleepingPlaces" bson:"sleepingPlaces"`
		AllowedNumberOfPeople int                 `json:"allowedNumberOfPeople" bson:"allowedNumberOfPeople"`
		AllowedNumberOfPets   int                 `json:"allowedNumberOfPets" bson:"allowedNumberOfPets"`
		TopAttributes         []AttributeReadOnly `json:"topAttributes" bson:"topAttributes"`
		Attributes            []AttributeReadOnly `json:"attributes" bson:"attributes"`
	}
)

func NewAttributeCollection() *AttributeCollection {
	attribute := AttributeCollection{
		RoomSize:              0,
		SleepingPlaces:        0,
		AllowedNumberOfPeople: 0,
		Bathrooms:             0,
		AllowedNumberOfPets:   0,
		TopAttributeIDs:       []string{},
		AttributeIDs:          []string{},
	}

	return &attribute
}

func (ac *AttributeCollection) UpdateAttributeIDs(attributeIDs ...string) {
	updatedTopAttributeIDs := []string{}

	// we need to ensure to update the entries from the top attributes as well
	for _, attributeID := range attributeIDs {
		for _, topAttributeID := range ac.TopAttributeIDs {
			shouldNotBeRemoved := attributeID == topAttributeID
			if shouldNotBeRemoved {
				updatedTopAttributeIDs = append(updatedTopAttributeIDs, topAttributeID)
			}
		}
	}

	ac.TopAttributeIDs = updatedTopAttributeIDs
	ac.AttributeIDs = attributeIDs
}

func (ac *AttributeCollection) UpdateTopAttributeIDs(topAttributeIDs ...string) {
	ac.TopAttributeIDs = topAttributeIDs
}
