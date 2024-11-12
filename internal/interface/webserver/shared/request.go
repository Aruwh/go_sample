package shared

import (
	"fewoserv/internal/infrastructure/common"
	"time"
)

type (
	DeleteRequest struct {
		ID string `json:"id"`
	}

	Translation struct {
		De_DE *string `json:"deDE"`
		En_GB *string `json:"enGB"`
		Fr_FR *string `json:"frFR"`
		It_IT *string `json:"itIT"`
	}

	Picture struct {
		ID          string       `json:"id"`
		Description *Translation `json:"description"`
		Raw         *string      `json:"raw"`
	}

	SaisonEntry struct {
		Type     int       `json:"type"`
		FromDate time.Time `json:"fromDate"`
		ToDate   time.Time `json:"toDate"`
	}

	Address struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		StreetName   string `json:"streetName"`
		StreetNumber string `json:"streetNumber"`
		Zip          string `json:"zip"`
		City         string `json:"city"`
		Country      string `json:"country"`
	}

	PhoneNumber struct {
		Type   *common.PhoneNumberType `json:"type" bson:"type"`
		Number *string                 `json:"number" bson:"number"`
	}
)
