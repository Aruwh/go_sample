package application

import (
	"fewoserv/internal/infrastructure/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	firstName     = "Horst"
	lastName      = "Schlemmer"
	feDestination = "T0pS€cReT"
	token         = "T0pS€cReT"
	locale        = common.DeDE
)

func TestNew(t *testing.T) {
	invitation := BuildEmailInvitationTemplate(locale, firstName, lastName, feDestination, token)

	assert.NotNil(t, invitation)
}
