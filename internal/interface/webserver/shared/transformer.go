package shared

import (
	"fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
)

func TransformRequestSaisonEntries(saisonEntries *[]SaisonEntry) *[]shared.SaisonEntry {
	isSomethingToDo := saisonEntries != nil
	if !isSomethingToDo {
		return nil
	}

	transformedSaisonEntries := []shared.SaisonEntry{}
	for _, saisonEntry := range *saisonEntries {
		transformedSaisonEntry := shared.SaisonEntry{
			Type:     common.SaisonType(saisonEntry.Type),
			FromDate: saisonEntry.FromDate,
			ToDate:   saisonEntry.ToDate,
		}

		transformedSaisonEntries = append(transformedSaisonEntries, transformedSaisonEntry)
	}

	return &transformedSaisonEntries
}

func TransformRequestTranslation(translation *Translation) *shared.Translation {
	isSomethingToDo := translation != nil
	if !isSomethingToDo {
		return nil
	}

	transformedPasswordUpdate := shared.Translation{
		De_DE: translation.De_DE,
		En_GB: translation.En_GB,
		Fr_FR: translation.Fr_FR,
		It_IT: translation.It_IT,
	}

	return &transformedPasswordUpdate
}

func TransformRequestPicture(picture *shared.Picture, base64RawData *string) *Picture {
	isSomethingToDo := picture != nil
	if !isSomethingToDo {
		return nil
	}

	transformedDescription := Translation{
		De_DE: picture.Description.De_DE,
		En_GB: picture.Description.En_GB,
		Fr_FR: picture.Description.Fr_FR,
		It_IT: picture.Description.It_IT,
	}

	transformedPictureUpdate := Picture{
		ID:          picture.ID,
		Description: &transformedDescription,
		Raw:         base64RawData,
	}

	return &transformedPictureUpdate
}
