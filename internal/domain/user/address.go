package user

import "fewoserv/pkg/mongodb"

type (
	Address struct {
		ID           string `json:"id" bson:"_id"`
		FirstName    string `json:"firstName" bson:"firstName"`
		LastName     string `json:"lastName" bson:"lastName"`
		StreetName   string `json:"streetName" bson:"streetName"`
		StreetNumber string `json:"streetNumber" bson:"streetNumber"`
		Zip          string `json:"zip" bson:"zip"`
		City         string `json:"city" bson:"city"`
		Country      string `json:"country" bson:"country"`
	}
)

func NewAddress(firstName, lastName, streetName, streetNumber, zip, city, country string) *Address {
	address := Address{
		ID:           mongodb.NewID(),
		FirstName:    firstName,
		LastName:     lastName,
		StreetName:   streetName,
		StreetNumber: streetNumber,
		Zip:          zip,
		City:         city,
		Country:      country,
	}

	return &address
}

func (a *Address) UpdateAddress(address *Address) {
	shouldBeUpdated := address.FirstName != "" && a.FirstName != address.FirstName
	if shouldBeUpdated {
		a.FirstName = address.FirstName
	}

	shouldBeUpdated = address.LastName != "" && a.LastName != address.LastName
	if shouldBeUpdated {
		a.LastName = address.LastName
	}

	shouldBeUpdated = address.StreetName != "" && a.StreetName != address.StreetName
	if shouldBeUpdated {
		a.StreetName = address.StreetName
	}

	shouldBeUpdated = address.StreetNumber != "" && a.StreetNumber != address.StreetNumber
	if shouldBeUpdated {
		a.StreetNumber = address.StreetNumber
	}

	shouldBeUpdated = address.Zip != "" && a.Zip != address.Zip
	if shouldBeUpdated {
		a.Zip = address.Zip
	}

	shouldBeUpdated = address.City != "" && a.City != address.City
	if shouldBeUpdated {
		a.City = address.City
	}

	shouldBeUpdated = address.Country != "" && a.Country != address.Country
	if shouldBeUpdated {
		a.Country = address.Country
	}
}
