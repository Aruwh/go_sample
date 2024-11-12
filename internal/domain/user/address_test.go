package user_test

// import (
// 	"fewoserv/internal/domain/user"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	streetName   = "Main Street"
// 	streetNumber = "31c"
// 	zip          = "1M1"
// 	city         = "Douglas"
// 	country      = "Isle of man"
// )

// func TestNewAddress(t *testing.T) {
// 	testNewAddress := user.NewAddress(firstName, lastName, streetName, streetNumber, zip, city, country)

// 	assert.NotNil(t, testNewAddress)
// 	assert.Equal(t, testNewAddress.FirstName, firstName)
// 	assert.Equal(t, testNewAddress.LastName, lastName)
// 	assert.Equal(t, testNewAddress.StreetName, streetName)
// 	assert.Equal(t, testNewAddress.StreetNumber, streetNumber)
// 	assert.Equal(t, testNewAddress.Zip, zip)
// 	assert.Equal(t, testNewAddress.City, city)
// 	assert.Equal(t, testNewAddress.Country, country)
// }

// func TestUpdateAddress(t *testing.T) {
// 	testNewAddress := user.NewAddress(firstName, lastName, streetName, streetNumber, zip, city, country)

// 	updateAddress := user.Address{
// 		FirstName: "Albert",
// 		LastName:  "Einstein",
// 	}

// 	testNewAddress.UpdateAddress(&updateAddress)
// 	assert.Equal(t, testNewAddress.FirstName, updateAddress.FirstName)
// 	assert.Equal(t, testNewAddress.LastName, updateAddress.LastName)
// 	assert.Equal(t, testNewAddress.StreetName, streetName)
// 	assert.Equal(t, testNewAddress.StreetNumber, streetNumber)
// 	assert.Equal(t, testNewAddress.Zip, zip)
// 	assert.Equal(t, testNewAddress.City, city)
// 	assert.Equal(t, testNewAddress.Country, country)
// }
