package booking

import (
	"time"
)

type (
	Message struct {
		Timestamp time.Time `json:"timestamp" bson:"timestamp"`
		Text      string    `json:"text" bson:"text"`
	}
)
