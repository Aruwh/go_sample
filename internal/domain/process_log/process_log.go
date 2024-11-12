package processlog

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"
)

const ()

type (
	ProcessLog struct {
		ID       string           `json:"id" bson:"_id"`
		UserID   string           `json:"userID" bson:"userID"`
		RecordID *string          `json:"recordID" bson:"recordID"`
		Action   common.Action    `json:"action" bson:"action"`
		Domain   common.Domain    `json:"domain" bson:"domain"`
		Value    string           `json:"value" bson:"value"`
		Created  shared.TimeStamp `json:"created" bson:"created"`
	}
)

func New(userID, value string, action common.Action, domain common.Domain, recordID *string) *ProcessLog {
	timestamp := shared.NewTimeStamp(&userID)

	log := ProcessLog{
		ID:       mongodb.NewID(),
		UserID:   userID,
		RecordID: recordID,
		Action:   action,
		Domain:   domain,
		Value:    value,
		Created:  timestamp,
	}

	return &log
}
