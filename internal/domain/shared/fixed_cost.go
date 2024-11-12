package shared

type (
	FixedCostEntry struct {
		ID    string      `json:"id" bson:"_id"`
		Name  Translation `json:"name" bson:"name"`
		Price float32     `json:"price" bson:"price"`
	}

	FixCost struct {
		CleaningFee FixedCostEntry `json:"cleaningFee" bson:"cleaningFee"`
		Total       float32        `json:"total" bson:"total"`
	}
)

func buildCleaningFee() FixedCostEntry {
	var (
		deDE = "Reinigungspauschale"
		enGB = "Cleaning fee"
		frFR = "Frais de nettoyage"
		itIT = "Spese di pulizia"
	)
	return FixedCostEntry{
		ID:    "0",
		Name:  Translation{De_DE: &deDE, En_GB: &enGB, Fr_FR: &frFR, It_IT: &itIT},
		Price: 0,
	}
}

func NewFixCost() FixCost {
	fixCost := FixCost{CleaningFee: buildCleaningFee(), Total: 0}

	return fixCost
}

func (fc *FixCost) CalculateTotal() float32 {
	var total float32 = 0
	total += fc.CleaningFee.Price

	fc.Total = total

	return total
}
