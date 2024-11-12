package user

import (
	userDomain "fewoserv/internal/domain/user"
	"fewoserv/internal/interface/webserver/shared"
)

func TransformRequestAddress(requestAddress shared.Address) *userDomain.Address {
	transformedAddress := userDomain.NewAddress(requestAddress.FirstName, requestAddress.LastName, requestAddress.StreetName, requestAddress.StreetNumber, requestAddress.Zip, requestAddress.City, requestAddress.Country)
	return transformedAddress
}

func TransformRequestPhoneNumbe(requestPhoneNumber shared.PhoneNumber) *userDomain.PhoneNumber {
	transformedPhoneNumber := userDomain.NewPhoneNumber(*requestPhoneNumber.Type, *requestPhoneNumber.Number)
	return transformedPhoneNumber
}

func TransformRequestAddresses(requestAddresses []shared.Address) []userDomain.Address {
	transformedAdresses := []userDomain.Address{}

	for _, requestAddress := range requestAddresses {
		transformedAdress := TransformRequestAddress(requestAddress)
		transformedAdresses = append(transformedAdresses, *transformedAdress)
	}

	return transformedAdresses
}

func TransformRequestPhoneNumbers(requestPhoneNumbers []shared.PhoneNumber) []userDomain.PhoneNumber {
	transformedPhoneNumbers := []userDomain.PhoneNumber{}

	for _, requestPhoneNumber := range requestPhoneNumbers {
		transformedPhoneNumber := TransformRequestPhoneNumbe(requestPhoneNumber)
		transformedPhoneNumbers = append(transformedPhoneNumbers, *transformedPhoneNumber)
	}

	return transformedPhoneNumbers
}
