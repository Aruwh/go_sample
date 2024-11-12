package realEstate_test

import (
	realEstate "fewoserv/internal/domain/real_estate"
	"fewoserv/internal/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	creatorID = "öaskdöfksdf"
	pictureID = "234234234köl"

	realEstateName = "FirstRealEstate"

	realEstateDescription = "I am the first real estate ..."
	description           = shared.NewTranslation(realEstateDescription)
)

func TestNew(t *testing.T) {
	testRealEstate, err := realEstate.New(creatorID, &realEstateName, &pictureID, description)

	assert.NoError(t, err)
	assert.NotNil(t, testRealEstate)

	assert.NotNil(t, testRealEstate.Name)
	assert.Equal(t, *testRealEstate.Name, realEstateName)

	assert.NotNil(t, testRealEstate.Description)
	assert.Equal(t, *testRealEstate.Description.De_DE, realEstateDescription)
	assert.Equal(t, *testRealEstate.Description.Fr_FR, realEstateDescription)
	assert.Equal(t, *testRealEstate.Description.En_GB, realEstateDescription)
	assert.Equal(t, *testRealEstate.Description.It_IT, realEstateDescription)
}

func TestNew_ShouldReturnErrorBcsNoNamePass(t *testing.T) {
	testRealEstate, err := realEstate.New(creatorID, nil, nil, description)
	assert.ErrorIs(t, err, realEstate.ErrRealEstateNoName)
	assert.Nil(t, testRealEstate)
}

func TestNew_ShouldReturnErrorBcsNoDescriptionPass(t *testing.T) {
	testRealEstate, err := realEstate.New(creatorID, &realEstateName, nil, nil)
	assert.ErrorIs(t, err, realEstate.ErrRealEstateNoDescription)
	assert.Nil(t, testRealEstate)
}

func TestUpdate(t *testing.T) {
	testRealEstate, _ := realEstate.New(creatorID, &realEstateName, &pictureID, description)

	testUpdateName := "I am the AlpenHooooooorn"

	testPictureID := "asdlfjlsafjlasjkdfsa"

	testRealEstate.Update(&testPictureID, &testUpdateName, nil)

	assert.Equal(t, *testRealEstate.Name, testUpdateName)
	assert.Equal(t, *testRealEstate.PictureID, testPictureID)
}
