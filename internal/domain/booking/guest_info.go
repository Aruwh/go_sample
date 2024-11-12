package booking

type (
	GuestInfo struct {
		AdultAmount int `json:"adultAmount" bson:"adultAmount"`
		ChildAmount int `json:"childAmount" bson:"childAmount"`
		PetAmount   int `json:"petAmount" bson:"petAmount"`
	}
)

func NewGuestInfo(adultAmount, childAmount, petAmount *int) *GuestInfo {
	var usedAdultAmount, usedChildAmount, usedPetAmount int = 0, 0, 0

	if adultAmount != nil {
		usedAdultAmount = *adultAmount
	}

	if childAmount != nil {
		usedChildAmount = *childAmount
	}

	if petAmount != nil {
		usedPetAmount = *petAmount
	}

	guestInfo := GuestInfo{
		AdultAmount: usedAdultAmount,
		ChildAmount: usedChildAmount,
		PetAmount:   usedPetAmount,
	}

	return &guestInfo
}

func (gf *GuestInfo) Update(adultAmount, childAmount, petAmount *int) {
	if adultAmount != nil {
		gf.AdultAmount = *adultAmount
	}

	if childAmount != nil {
		gf.ChildAmount = *childAmount
	}

	if petAmount != nil {
		gf.PetAmount = *petAmount
	}
}
