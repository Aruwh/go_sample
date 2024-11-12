package booking

import (
	"fewoserv/internal/domain/apartment"
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"time"
)

type (
	HistoryItem struct {
		SaisonType common.SaisonType `json:"saisonType"`
		Price      float64           `json:"price"`
		Tax        float64           `json:"tax"`
	}

	PriceSummary struct {
		Total    float64        `json:"total" bson:"total"`
		Vat      int8           `json:"vat" bson:"vat"`
		Tax      float64        `json:"tax" bson:"tax"`
		FixCosts shared.FixCost `json:"fixCost" bson:"fixCosts"`
		Discount float64        `json:"discount" bson:"discount"`
		History  []HistoryItem  `json:"history" bson:"history"`
	}
)

func NewPriceSummary() *PriceSummary {
	priceSummary := PriceSummary{
		Total:    0,
		Vat:      0,
		Tax:      0,
		Discount: 0,
	}

	return &priceSummary
}

func (ps *PriceSummary) CalculatePriceSummary(fromDate, toDate time.Time, saisonPrice *apartment.SaisonPrice, datesWithSeasonTypes map[time.Time]common.SaisonType) *PriceSummary {
	var (
		total        float64       = 0
		tax          float64       = 0
		historyItems []HistoryItem = []HistoryItem{}
	)

	currentDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 6, 0, 0, 0, time.UTC)

	// exeption when only one day was booked
	if fromDate.Year() == toDate.Year() && fromDate.Month() == toDate.Month() && fromDate.Day() == toDate.Day() {
		saisonType, isSeasonAvailable := datesWithSeasonTypes[currentDate]
		bruttoPriceInfo := saisonPrice.GetBruttoPriceInfo(saisonType, isSeasonAvailable)

		total += bruttoPriceInfo.Price
		tax += bruttoPriceInfo.Tax

		historyItem := HistoryItem{SaisonType: saisonType, Price: bruttoPriceInfo.Price, Tax: bruttoPriceInfo.Tax}
		historyItems = append(historyItems, historyItem)
	}

	for !currentDate.After(toDate) {
		saisonType, isSeasonAvailable := datesWithSeasonTypes[currentDate]
		bruttoPriceInfo := saisonPrice.GetBruttoPriceInfo(saisonType, isSeasonAvailable)

		total += bruttoPriceInfo.Price
		tax += bruttoPriceInfo.Tax

		historyItem := HistoryItem{SaisonType: saisonType, Price: bruttoPriceInfo.Price, Tax: bruttoPriceInfo.Tax}
		historyItems = append(historyItems, historyItem)

		// jump to next date
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	saisonPrice.FixCost.CalculateTotal()

	// we add the fix cost on the end of the calculation
	ps.Total = total + float64(saisonPrice.FixCost.Total)

	ps.Tax = tax
	ps.Vat = saisonPrice.VatPercent
	ps.FixCosts = saisonPrice.FixCost
	ps.History = historyItems

	return ps
}
