package attribute

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/pkg/mongodb"
)

type (
	Attribute struct {
		ID      string              `json:"id" bson:"_id"`
		Name    *shared.Translation `json:"name" bson:"name"`
		Svg     *string             `json:"svg" bson:"svg"`
		Created *shared.TimeStamp   `json:"created" bson:"created"`
		Edited  *shared.TimeStamp   `json:"edited" bson:"edited"`
	}
)

func New(createrID string, name *shared.Translation, svg *string) *Attribute {
	timestamp := shared.NewTimeStamp(&createrID)

	attribute := Attribute{
		ID:      mongodb.NewID(),
		Name:    name,
		Svg:     svg,
		Created: &timestamp,
		Edited:  &timestamp,
	}

	return &attribute
}

func (a *Attribute) Update(name *shared.Translation, svg *string) {
	a.Name.Update(name)

	shouldBeUpdated := svg != nil
	if shouldBeUpdated {
		a.Svg = svg
	}
}
