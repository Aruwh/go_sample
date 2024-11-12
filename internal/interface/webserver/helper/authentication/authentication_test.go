package helper_test

import (
	"fewoserv/internal/infrastructure/common"
	helper "fewoserv/internal/interface/webserver/helper/authentication"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	userID                  = "234köasok233"
	permissions             = []common.RequestPermission{common.BOOKING_VIEW, common.SETTINGS_VIEW}
	expirationTimeInSeconds = 1
)

func TestGenerateJwt(t *testing.T) {
	token, err := helper.GenerateJwt(&userID, common.ADMINISTRATOR, permissions, common.ItIT, &expirationTimeInSeconds)

	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestNewTokenWithErrorBecauseTokenIsExpired(t *testing.T) {
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc4Nzc4MzksImlhdCI6MTY5Nzg3NzgzOCwiaXNzIjoiRmVXb1NlcnYuY2giLCJwZXJtaXNzaW9ucyI6WyJBRE1JTiIsIk9XTkVSIl0sInN1YiI6IjIzNGvDtmFzb2syMzMifQ.NGLuWiXlDc4RafCX2NOs7EeY7O51Rxz0ZTddQLqk2EM"
	token, err := helper.ValidateAndTransformToken(testToken)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.ErrorAs(t, err, &helper.ErrNotAuthorised)
}

func TestNewTokenWithErrorBecauseTokenSignatureWasManipulated(t *testing.T) {
	testToken := "fyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc4Nzc4MzksImlhdCI6MTY5Nzg3NzgzOCwiaXNzIjoiRmVXb1NlcnYuY2giLCJwZXJtaXNzaW9ucyI6WyJBRE1JTiIsIk9XTkVSIl0sInN1YiI6IjIzNGvDtmFzb2syMzMifQ.NGLuWiXlDc4RafCX2NOs7EeY7O51Rxz0ZTddQLqk2EM"
	token, err := helper.ValidateAndTransformToken(testToken)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.ErrorAs(t, err, &helper.ErrMalformedToken)
}

func TestNewToken(t *testing.T) {
	// Note: infinity token
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTgyNTcyOTMsImlhdCI6MTY5Nzg4MDg3MywiaXNzIjoiRmVXb1NlcnYuY2giLCJwZXJtaXNzaW9ucyI6WyJBRE1JTiIsIk9XTkVSIl0sInN1YiI6IjIzNGvDtmFzb2syMzMifQ.rWHj6YN_pqK8CvS_TO84r8JdNNeKzrsQAR9Bdfthq0A"
	token, err := helper.ValidateAndTransformToken(testToken)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)
}
