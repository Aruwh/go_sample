package permission

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/pkg/mongodb"
)

type (
	Permission struct {
		ID          string             `json:"id" bson:"_id"`
		Name        *string            `json:"name" bson:"name"`
		Description shared.Translation `json:"description" bson:"description"`
	}
)

func New(name string, description shared.Translation) *Permission {
	permission := Permission{
		ID:          mongodb.NewID(),
		Name:        &name,
		Description: description,
	}

	return &permission
}

func (p *Permission) Update(permission *Permission) {
	shouldBeUpdated := permission.Name != nil && permission.Name != p.Name
	if shouldBeUpdated {
		p.Name = permission.Name
	}

	p.Description.Update(&permission.Description)
}
