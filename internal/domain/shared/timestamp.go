package shared

import "time"

type (
	TimeStamp struct {
		UserID string    `json:"userID" bson:"userID"`
		Time   time.Time `json:"time" bson:"time"`
	}
)

func NewTimeStamp(userID *string) TimeStamp {
	var timeStamp = TimeStamp{}
	timeStamp.Time = time.Now()

	if userID == nil {
		timeStamp.UserID = DEFAULT_USER_ID
	} else {
		timeStamp.UserID = *userID
	}

	return timeStamp
}

func (ts *TimeStamp) Update(updaterID string) {
	ts.UserID = updaterID
	ts.Time = time.Now()
}
