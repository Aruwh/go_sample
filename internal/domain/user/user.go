package user

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/pkg/mongodb"
	"time"
)

type (
	User struct {
		ID           string           `json:"id" bson:"_id"`
		IsActive     *bool            `json:"isActive" bson:"isActive"`
		Email        *string          `json:"email" bson:"email"`
		Password     *string          `json:"password" bson:"password"`
		Sex          *common.Sex      `json:"sex" bson:"sex"`
		BirthDate    *time.Time       `json:"birthDate" bson:"birthDate"`
		Locale       *common.Locale   `json:"locale" bson:"locale"`
		FirstName    *string          `json:"firstName" bson:"firstName"`
		LastName     *string          `json:"lastName" bson:"lastName"`
		Addresses    []Address        `json:"addresses" bson:"addresses"`
		PhoneNumbers []PhoneNumber    `json:"phoneNumbers" bson:"phoneNumbers"`
		Created      shared.TimeStamp `json:"created" bson:"created"`
		Edited       shared.TimeStamp `json:"edited" bson:"edited"`
	}
)

func New(email, password, firstName, lastName string, birthDate time.Time, locale common.Locale, sex common.Sex) *User {
	userID := mongodb.NewID()

	timestamp := shared.NewTimeStamp(&userID)
	user := User{
		ID:           mongodb.NewID(),
		Email:        &email,
		Password:     &password,
		Sex:          &sex,
		BirthDate:    &birthDate,
		Locale:       &locale,
		FirstName:    &firstName,
		LastName:     &lastName,
		Addresses:    []Address{},
		PhoneNumbers: []PhoneNumber{},
		Created:      timestamp,
		Edited:       timestamp,
	}

	return &user
}

func (u *User) UpdatePassword(oldPassword, newPassword, repeatedNewPassword string) error {
	decryptedOldPassword, err := utils.DecryptPwd(oldPassword)
	if err != nil {
		return err
	}

	isPasswordResetAllowed := u.Password == decryptedOldPassword
	if !isPasswordResetAllowed {
		return ErrUserPasswordMissMatch
	}

	areNewPasswordsEqual := newPassword == repeatedNewPassword
	if !areNewPasswordsEqual {
		return ErrUserNewPasswordsAreNotEqual
	}

	decryptedNewPassword, err := utils.DecryptPwd(newPassword)
	if err != nil {
		return err
	}

	u.Password = decryptedNewPassword

	return nil
}

func (u *User) Update(user *User) {
	shouldBeUpdated := user.Sex != nil && u.Sex != user.Sex
	if shouldBeUpdated {
		u.Sex = user.Sex
	}

	shouldBeUpdated = user.BirthDate != nil && u.BirthDate != user.BirthDate
	if shouldBeUpdated {
		u.BirthDate = user.BirthDate
	}

	shouldBeUpdated = user.Locale != nil && u.Locale != user.Locale
	if shouldBeUpdated {
		u.Locale = user.Locale
	}

	shouldBeUpdated = u.FirstName != nil && u.FirstName == user.FirstName
	if shouldBeUpdated {
		u.FirstName = user.FirstName
	}

	shouldBeUpdated = u.LastName != nil && u.LastName == user.LastName
	if shouldBeUpdated {
		u.LastName = user.LastName
	}
}

func (u *User) AddAddress(address *Address) {
	u.Addresses = append(u.Addresses, *address)
}

func (u *User) RemoveAddress(addressIDs ...string) {
	reducedAddresses := []Address{}

	for _, addressID := range addressIDs {
		for _, address := range u.Addresses {
			shouldNotBeRemoved := address.ID != addressID
			if shouldNotBeRemoved {
				reducedAddresses = append(reducedAddresses, address)
			}
		}
	}

	u.Addresses = reducedAddresses
}

func (u *User) AddPhoneNumber(phoneNumber *PhoneNumber) {
	u.PhoneNumbers = append(u.PhoneNumbers, *phoneNumber)
}

func (u *User) RemovePhoneNumber(phoneNumberIDs ...string) {
	reducedPhoneNumber := []PhoneNumber{}

	for _, phoneNumberID := range phoneNumberIDs {
		for _, phoneNumber := range u.PhoneNumbers {
			shouldNotBeRemoved := phoneNumber.ID != phoneNumberID
			if shouldNotBeRemoved {
				reducedPhoneNumber = append(reducedPhoneNumber, phoneNumber)
			}
		}
	}

	u.PhoneNumbers = reducedPhoneNumber
}
