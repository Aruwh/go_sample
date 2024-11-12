package picturecache

import (
	"fewoserv/internal/domain/shared"
	"fmt"
)

type (
	PictureCacheItem struct {
		ID          string                    `json:"id"`
		Description shared.Translation        `json:"description"`
		Destination shared.PictureDestination `json:"-"`
		Raw         string                    `json:"raw"`
	}

	IPictureCache interface {
		Add(variant string, picture *shared.Picture)
		Get(variant, recordID string) *shared.Picture
	}

	PictureCache struct {
		Entries map[string]*PictureCacheItem
	}
)

func NewCacheItem(picture *shared.Picture) *PictureCacheItem {
	cacheItem := PictureCacheItem{
		ID:          picture.ID,
		Description: *picture.Description,
		Destination: *picture.Destination,
		Raw:         *picture.Raw,
	}

	return &cacheItem
}

func CacheItemToPicture(cacheItem PictureCacheItem) *shared.Picture {
	picture := shared.Picture{
		ID:          cacheItem.ID,
		Description: &cacheItem.Description,
		Destination: &cacheItem.Destination,
		Raw:         &cacheItem.Raw,
	}

	return &picture
}

func New() IPictureCache {
	pictureCache := PictureCache{Entries: make(map[string]*PictureCacheItem)}

	return &pictureCache
}

func buildIdentifier(variant, recordID string) string {
	return fmt.Sprintf("%s%s", variant, recordID)
}

func (pc *PictureCache) Add(variant string, picture *shared.Picture) {
	identifier := buildIdentifier(variant, picture.ID)
	pc.Entries[identifier] = NewCacheItem(picture)
}

func (pc *PictureCache) Get(variant, recordID string) *shared.Picture {
	identifier := buildIdentifier(variant, recordID)
	cachedItem := pc.Entries[identifier]

	if cachedItem == nil {
		return nil
	}

	return CacheItemToPicture(*cachedItem)
}
