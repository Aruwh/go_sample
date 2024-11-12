package news

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/pkg/mongodb"
	"time"
)

type (
	News struct {
		ID        string              `json:"id" bson:"_id"`
		Title     *shared.Translation `json:"title" bson:"title"`
		Content   *shared.Translation `json:"content" bson:"content"`
		PublishAt *time.Time          `json:"publishAt" bson:"publishAt"`
		Active    *bool               `json:"active" bson:"active"`
		Created   *shared.TimeStamp   `json:"created" bson:"created"`
		Edited    *shared.TimeStamp   `json:"edited" bson:"edited"`
	}
)

func New(createrID string, title, content shared.Translation, publishAt time.Time, active bool) *News {
	timestamp := shared.NewTimeStamp(&createrID)

	post := News{
		ID:        mongodb.NewID(),
		Title:     &title,
		Content:   &content,
		PublishAt: &publishAt,
		Active:    &active,
		Created:   &timestamp,
		Edited:    &timestamp,
	}

	return &post
}

func (p *News) Update(title, content *shared.Translation, publishAt *time.Time, active *bool) {
	p.Title.Update(title)
	p.Content.Update(content)

	shouldBeUpdated := active != nil && *p.Active != *active
	if shouldBeUpdated {
		p.Active = active
	}

	shouldBeUpdated = publishAt != nil
	if shouldBeUpdated {
		p.PublishAt = publishAt
	}
}
