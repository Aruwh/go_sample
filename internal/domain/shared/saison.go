package shared

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"
	"fmt"
	"time"
)

var logApartment = logger.New("SAISON")

type (
	SaisonEntry struct {
		Type     common.SaisonType `json:"type" bson:"type"`
		FromDate time.Time         `json:"fromDate" bson:"fromDate"`
		ToDate   time.Time         `json:"toDate" bson:"toDate"`
	}

	Saison struct {
		ID      string        `json:"id" bson:"_id"`
		Year    int           `json:"year" bson:"year"`
		Entries []SaisonEntry `json:"entries" bson:"entries"`
		Created *TimeStamp    `json:"created" bson:"created"`
		Edited  *TimeStamp    `json:"edited" bson:"edited"`
	}
)

func NewSaison(createrID string, year int, entries []SaisonEntry) *Saison {
	recordID := mongodb.NewID()
	timeStamp := NewTimeStamp(&createrID)

	var saison = Saison{
		ID:      recordID,
		Year:    year,
		Entries: entries,
		Created: &timeStamp,
		Edited:  &timeStamp,
	}

	logApartment.Debug(fmt.Sprintf("Created new saison %v", saison))

	return &saison
}

func (s *Saison) Update(year *int, entries *[]SaisonEntry) {
	shouldBeUpdated := year != nil && s.Year != *year
	if shouldBeUpdated {
		s.Year = *year
	}

	shouldBeUpdated = entries != nil
	if shouldBeUpdated {
		s.Entries = *entries
	}
}
