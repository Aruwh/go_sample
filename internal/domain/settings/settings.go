package settings

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/pkg/mongodb"
)

type (
	Settings struct {
		ID                  string              `json:"id" bson:"_id"`
		StornoDayRange      int8                `json:"stornoDayRange" bson:"stornoDayRange"`
		BookingNumber       int                 `json:"bookingNumber" bson:"bookingNumber"`
		NotificationMessage *shared.Translation `json:"notificationMessage" bson:"notificationMessage"`
	}
)

func New() *Settings {
	siteSettings := Settings{
		ID:                  mongodb.NewID(),
		StornoDayRange:      0,
		BookingNumber:       0,
		NotificationMessage: nil,
	}

	return &siteSettings
}

func (ss *Settings) UpdateStornoDayRange(stornoDayRange int8) {
	if stornoDayRange != ss.StornoDayRange {
		ss.StornoDayRange = stornoDayRange
	}
}

func (ss *Settings) UpdateNotificationMessage(message *shared.Translation) {
	ss.NotificationMessage.Update(message)
}
