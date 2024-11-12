package helper

import (
	"fewoserv/internal/infrastructure/common"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Identity struct {
		Permissions []common.RequestPermission
		UserID      string
		Type        common.AudianceType
	}
)

func TransformTokenToIdentity(token *jwt.Token) *Identity {
	mapWithClaims := token.Claims.(jwt.MapClaims)

	userID := mapWithClaims["sub"].(string)

	permissions := []common.RequestPermission{}
	for _, permission := range mapWithClaims["permissions"].([]interface{}) {
		permissions = append(permissions, common.RequestPermission(permission.(string)))
	}

	identity := Identity{
		Permissions: permissions,
		UserID:      userID,
		Type:        common.AudianceType(mapWithClaims["aud"].(string)),
	}

	return &identity
}

func (i *Identity) EnsureRequestPermission(permissions ...common.RequestPermission) error {
	var (
		err                      error = nil
		permissionMatchedCounter       = 0
	)

	// if the current user has admin permissions, the check can pass immediatly
	isSuperAdmin := i.Type == common.AUDIANCE_SUPER_ADMIN_USER
	if isSuperAdmin {
		return nil
	}

	for _, identityPermission := range i.Permissions {
		for _, permission := range permissions {

			hasMatched := permission == identityPermission
			if hasMatched {
				permissionMatchedCounter++
			}
		}
	}

	permissionAllowed := len(permissions) == permissionMatchedCounter
	if !permissionAllowed {
		err = ErrNoPermission
	}

	return err
}
