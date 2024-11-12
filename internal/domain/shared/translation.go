package shared

type (
	Translation struct {
		De_DE *string `json:"deDE" bson:"deDE"`
		En_GB *string `json:"enGB" bson:"enGB"`
		Fr_FR *string `json:"frFR" bson:"frFR"`
		It_IT *string `json:"itIT" bson:"itIT"`
	}
)

// TODO: umbauen, sodass ne struct übergeben wird
func NewTranslation(text string) *Translation {
	translation := Translation{
		De_DE: &text,
		Fr_FR: &text,
		En_GB: &text,
		It_IT: &text,
	}

	return &translation
}

func (t *Translation) Update(translation *Translation) {
	canIDoSomething := translation != nil
	if !canIDoSomething {
		return
	}

	shouldBeUpdated := translation.De_DE != nil && translation.De_DE != t.De_DE
	if shouldBeUpdated {
		t.De_DE = translation.De_DE
	}

	shouldBeUpdated = translation.En_GB != nil && translation.En_GB != t.En_GB
	if shouldBeUpdated {
		t.En_GB = translation.En_GB
	}

	shouldBeUpdated = translation.Fr_FR != nil && translation.Fr_FR != t.Fr_FR
	if shouldBeUpdated {
		t.Fr_FR = translation.Fr_FR
	}

	shouldBeUpdated = translation.It_IT != nil && translation.It_IT != t.It_IT
	if shouldBeUpdated {
		t.It_IT = translation.It_IT
	}
}
