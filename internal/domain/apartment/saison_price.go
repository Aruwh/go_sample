package apartment

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"
	"fmt"
)

type (
	// All prices are netto
	SaisonPrice struct {
		ID            string         `json:"id" bson:"_id"`
		BasePrice     float64        `json:"basePrice" bson:"basePrice"`
		LowPrice      float64        `json:"lowPrice" bson:"lowPrice"`
		LowPercent    float64        `json:"lowPercent" bson:"lowPercent"`
		MiddlePrice   float64        `json:"middlePrice" bson:"middlePrice"`
		MiddlePercent float64        `json:"middlePercent" bson:"middlePercent"`
		HighPrice     float64        `json:"highPrice" bson:"highPrice"`
		HighPercent   float64        `json:"highPercent" bson:"highPercent"`
		PeakPrice     float64        `json:"peakPrice" bson:"peakPrice"`
		PeakPercent   float64        `json:"peakPercent" bson:"peakPercent"`
		Total         float64        `json:"total" bson:"total"`
		FixCost       shared.FixCost `json:"fixCost" bson:"fixCost"`
		VatPercent    int8           `json:"vatPercent" bson:"vatPercent"`
		VatPrice      int8           `json:"vatPrice" bson:"vatPrice"`
	}

	BruttoPriceInfo struct {
		Price   float64        `json:"price" bson:"price"`
		FixCost shared.FixCost `json:"fixCost" bson:"fixCost"`
		Tax     float64        `json:"tax" bson:"tax"`
		Vat     int8           `json:"vat" bson:"vat"`
	}
)

func calculateTotal(basePrice float64, fixCost shared.FixCost) float64 {
	var total float64 = basePrice

	total += float64(fixCost.CleaningFee.Price)

	return total
}

func NewSaisonPrice(basePrice, lowPrice, middlePrice, highPrice, peakPrice *float64, fixCost *shared.FixCost) *SaisonPrice {
	var (
		usedBasePrice   float64        = 0
		usedLowPrice    float64        = 0
		usedMiddlePrice float64        = 0
		usedHighPrice   float64        = 0
		usedPeakPrice   float64        = 0
		usedFixCost     shared.FixCost = shared.NewFixCost()
	)

	if lowPrice != nil {
		usedBasePrice = *basePrice
	}

	if lowPrice != nil {
		usedLowPrice = *lowPrice
	}

	if middlePrice != nil {
		usedMiddlePrice = *middlePrice
	}

	if highPrice != nil {
		usedHighPrice = *highPrice
	}

	if peakPrice != nil {
		usedPeakPrice = *peakPrice
	}

	if fixCost != nil {
		usedFixCost = *fixCost
	}

	total := calculateTotal(usedBasePrice, usedFixCost)

	saisonPrice := SaisonPrice{
		ID:          mongodb.NewID(),
		BasePrice:   usedBasePrice,
		LowPrice:    usedLowPrice,
		MiddlePrice: usedMiddlePrice,
		HighPrice:   usedHighPrice,
		PeakPrice:   usedPeakPrice,
		FixCost:     usedFixCost,
		Total:       total,
		VatPercent:  0,
	}

	return &saisonPrice
}

func (sp *SaisonPrice) UpdatePrice(saison common.SaisonType, price float64) error {
	switch saison {
	case common.SaisonLow:
		sp.LowPrice = price
	case common.SaisonMiddle:
		sp.MiddlePrice = price
	case common.SaisonHigh:
		sp.HighPrice = price
	case common.SaisonPeak:
		sp.PeakPrice = price
	default:
		return fmt.Errorf("%w: %v", ErrUnknownSaisonType, saison)
	}

	return nil
}

func (sp *SaisonPrice) UpdateVat(vat int8) {
	sp.VatPercent = vat
}

func (sp *SaisonPrice) GetBruttoPriceInfo(saison common.SaisonType, isSeasonAvailable bool) *BruttoPriceInfo {
	// if there is no season passed (isSeasonAvailable), we will fall back on the base price
	var nettoPrice float64 = sp.BasePrice

	if isSeasonAvailable {
		switch saison {
		case common.SaisonLow:
			nettoPrice = sp.LowPrice
		case common.SaisonMiddle:
			nettoPrice = sp.MiddlePrice
		case common.SaisonHigh:
			nettoPrice = sp.HighPrice
		case common.SaisonPeak:
			nettoPrice = sp.PeakPrice
		default:
			// TODO ERROR !!!
		}
	}

	tax := nettoPrice * 0.01 * float64(sp.VatPercent)

	bruttoPriceInfo := BruttoPriceInfo{
		Price:   nettoPrice + tax,
		Tax:     tax,
		FixCost: sp.FixCost,
		Vat:     sp.VatPercent,
	}

	return &bruttoPriceInfo
}
