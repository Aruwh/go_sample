package user

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"
)

type (
	PhoneNumber struct {
		ID     string                  `json:"id" bson:"_id"`
		Type   *common.PhoneNumberType `json:"type" bson:"type"`
		Number *string                 `json:"number" bson:"number"`
	}
)

func NewPhoneNumber(phoneNumberType common.PhoneNumberType, number string) *PhoneNumber {
	phoneNumber := PhoneNumber{
		ID:     mongodb.NewID(),
		Type:   &phoneNumberType,
		Number: &number,
	}

	return &phoneNumber
}

func (pn *PhoneNumber) UpdatePhoneNumber(phoneNumber *PhoneNumber) {
	shouldBeUpdated := phoneNumber.Type != nil && phoneNumber.Type != pn.Type
	if shouldBeUpdated {
		pn.Type = phoneNumber.Type
	}

	shouldBeUpdated = phoneNumber.Number != nil && phoneNumber.Number != pn.Number
	if shouldBeUpdated {
		pn.Number = phoneNumber.Number
	}
}
