package apartment_test

import (
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	creatorID     = "creatorID"
	realEstateID  = "realEstateID"
	ownerID       = "ownerID"
	apartmentName = "TestApartment"
)

func TestNew(t *testing.T) {
	testApartment := apartment.New(creatorID, realEstateID, ownerID, apartmentName)

	assert.NotNil(t, testApartment)
	assert.NotNil(t, testApartment.AttributeCollection)
	assert.NotNil(t, testApartment.Description)
	assert.NotNil(t, testApartment.PictureIDs)
	assert.Equal(t, *testApartment.Name, apartmentName)
	assert.Equal(t, *testApartment.Description.De_DE, "")
	assert.Equal(t, *testApartment.Description.Fr_FR, "")
	assert.Equal(t, *testApartment.Description.En_GB, "")
	assert.Equal(t, *testApartment.Description.It_IT, "")
	assert.Equal(t, testApartment.RealEstateID, realEstateID)
}

func TestUpdateRealEstateIDs(t *testing.T) {
	testApartment := apartment.New(creatorID, realEstateID, ownerID, apartmentName)
	assert.Equal(t, testApartment.RealEstateID, realEstateID)

	newRealEstateID := "newRealEstateID"
	testApartment.UpdateRealEstateID(newRealEstateID)

	assert.Equal(t, newRealEstateID, testApartment.RealEstateID)
}

func TestUpdateName(t *testing.T) {
	testApartment := apartment.New(creatorID, realEstateID, ownerID, apartmentName)

	testNewName := "Apartment Alpenblick"

	testApartment.UpdateName(&testNewName)
	assert.Equal(t, testNewName, testApartment.Name)
}

func TestUpdateDescription(t *testing.T) {
	testApartment := apartment.New(creatorID, realEstateID, ownerID, apartmentName)

	testNewDescriptionDe := "Bla bla DE"
	testNewDescriptionFr := "Blah Blah FR"
	testNewDescriptionEn := "Bluah Bluah EN"

	newDescription := shared.NewTranslation("")
	newDescription.De_DE = &testNewDescriptionDe
	newDescription.Fr_FR = &testNewDescriptionFr
	newDescription.En_GB = &testNewDescriptionEn

	testApartment.UpdateDescription(newDescription)
	assert.Equal(t, *newDescription.De_DE, testNewDescriptionDe)
	assert.Equal(t, *newDescription.Fr_FR, testNewDescriptionFr)
	assert.Equal(t, newDescription.It_IT, testApartment.Description.It_IT)
	assert.Equal(t, *newDescription.En_GB, testNewDescriptionEn)
}
