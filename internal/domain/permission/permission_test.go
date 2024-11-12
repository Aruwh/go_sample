package permission_test

import (
	"fewoserv/internal/domain/permission"
	"fewoserv/internal/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	name        = "BookingArea"
	description = shared.NewTranslation("Access to the booking area is granted to the user")
)

func TestNew(t *testing.T) {
	testPermission := permission.New(name, *description)

	assert.NotNil(t, testPermission)
	assert.Equal(t, &testPermission.Name, name)
	assert.Equal(t, &testPermission.Description, description)
}

func TestUpdate(t *testing.T) {
	testPermission := permission.New(name, *description)

	newName := "Booking"
	updatePermission := permission.Permission{
		Name: &newName,
	}

	testPermission.Update(&updatePermission)

	assert.Equal(t, testPermission.Name, updatePermission.Name)
}
