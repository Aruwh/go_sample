package shared_test

import (
	"fewoserv/internal/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTimestamp(t *testing.T) {
	timeStamp := shared.NewTimeStamp(nil)
	assert.NotNil(t, timeStamp)
	assert.Equal(t, timeStamp.UserID, "000000000000000000000000")
}
