package user_test

import (
	"fewoserv/internal/domain/user"
	"fewoserv/internal/infrastructure/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	number = "+41189302384"
)

func TestNewPhoneNumber(t *testing.T) {
	testPhoneNumber := user.NewPhoneNumber(common.Mobile, number)

	assert.NotNil(t, testPhoneNumber)
	assert.Equal(t, *testPhoneNumber.Number, number)
	assert.Equal(t, *testPhoneNumber.Type, common.Mobile)
}

func TestUpdatePhoneNumber(t *testing.T) {
	testPhoneNumber := user.NewPhoneNumber(common.Mobile, number)

	updateNumber := "09876543212"
	updatePhoneNumber := user.PhoneNumber{
		Number: &updateNumber,
	}

	testPhoneNumber.UpdatePhoneNumber(&updatePhoneNumber)
	assert.Equal(t, testPhoneNumber.Number, updatePhoneNumber.Number)
	assert.Equal(t, *testPhoneNumber.Type, common.Mobile)
}
