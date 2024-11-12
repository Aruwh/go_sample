package adminuser_test

import (
	adminuser "fewoserv/internal/domain/admin_user"
	"fewoserv/internal/infrastructure/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var creatorID = "aösdöjaksd"

func TestNewUser(t *testing.T) {
	email := "test@user.ch"
	password := "T0pS€cReT"
	firstName := "Ada"
	lastName := "Lovelace"

	user, err := adminuser.NewUser(creatorID, email, password, password, firstName, lastName, common.ADMINISTRATOR, common.FrFR)
	assert.NoError(t, err)

	assert.NotNil(t, user)
	assert.Equal(t, user.Email, email)
	assert.Equal(t, user.FirstName, firstName)
	assert.Equal(t, user.LastName, lastName)
}

func TestNewSuperUser(t *testing.T) {
	superUser, rawPassword, err := adminuser.NewSuperUser(&creatorID)
	assert.NoError(t, err)

	assert.NotNil(t, superUser)
	assert.NotNil(t, rawPassword)
	assert.Equal(t, superUser.Email, "fewoserv@admin.ch")
	assert.Equal(t, superUser.FirstName, "Super")
	assert.Equal(t, superUser.LastName, "User")
}

func TestAddPermissions(t *testing.T) {
	superUser, _, err := adminuser.NewSuperUser(&creatorID)
	assert.NoError(t, err)
	assert.Empty(t, superUser.Permissions)

	superUser.AddPermissions(common.ADMIN_USER_EDIT, common.BOOKING_VIEW)

	assert.Len(t, superUser.Permissions, 2)
}

func TestRemovePermissions(t *testing.T) {
	superUser, _, err := adminuser.NewSuperUser(&creatorID)
	assert.NoError(t, err)
	assert.Empty(t, superUser.Permissions)

	superUser.AddPermissions(common.ADMIN_USER_EDIT, common.BOOKING_VIEW)

	superUser.RemovePermissions(common.ADMIN_USER_EDIT)

	assert.Len(t, superUser.Permissions, 1)
	assert.Equal(t, superUser.Permissions[0], common.BOOKING_VIEW)
}

// TODO: add test with different pwds
