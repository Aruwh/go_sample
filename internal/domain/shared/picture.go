package shared

import (
	"encoding/base64"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/pkg/mongodb"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type (
	PictureDestinationEntry struct {
		Path string `json:"path" bson:"path"`
	}

	PictureDestination struct {
		Origin *PictureDestinationEntry `json:"origin" bson:"origin"`
		Large  *PictureDestinationEntry `json:"large" bson:"large"`
		Middle *PictureDestinationEntry `json:"middle" bson:"middle"`
		Small  *PictureDestinationEntry `json:"small" bson:"small"`
	}

	Picture struct {
		ID          string              `json:"id" bson:"_id"`
		Description *Translation        `json:"description" bson:"description"`
		Destination *PictureDestination `json:"-" bson:"destination"`
		Raw         *string             `json:"raw" bson:"-"`
	}

	FileInfo struct {
		ContentType string
		ModTime     time.Time
	}
)

func NewPicture(description *Translation) *Picture {
	usedDescription := *NewTranslation("")
	if description != nil {
		usedDescription = *description
	}

	newPicture := Picture{
		ID:          mongodb.NewID(),
		Description: &usedDescription,
		Destination: &PictureDestination{},
	}

	return &newPicture
}

func (p *PictureDestinationEntry) LoadFile() (*os.File, error) {
	imageFile, err := os.Open(p.Path)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()

	return imageFile, nil
}

func (p *PictureDestinationEntry) GetFileInfo() (*FileInfo, error) {
	// Stat the file to get its information
	fileInfo, err := os.Stat(p.Path)
	if err != nil {
		return nil, err
	}

	// get Content-Type based on the file extension
	contentType := http.DetectContentType([]byte(p.Path))

	return &FileInfo{ContentType: contentType, ModTime: fileInfo.ModTime()}, nil
}

func (p *PictureDestinationEntry) Tobase64() (*string, error) {
	imageFile, err := os.Open(p.Path)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()

	// Read the contents of the file
	content, err := io.ReadAll(imageFile)
	if err != nil {
		return nil, err
	}

	// Encode the content to Base64
	base64EncodedString := base64.StdEncoding.EncodeToString(content)
	htmlRdy := fmt.Sprintf("data:image/jpeg;base64,%s", base64EncodedString)

	return &htmlRdy, nil
}

func (p *Picture) Update(description *Translation) {
	p.Description.Update(description)
}

func (p *Picture) AddOriginDestination(path string) {
	p.Destination.Origin = &PictureDestinationEntry{Path: path}
}

func (p *Picture) AddSmallnDestination(path string) {
	p.Destination.Small = &PictureDestinationEntry{Path: path}
}

func (p *Picture) AddMiddleDestination(path string) {
	p.Destination.Middle = &PictureDestinationEntry{Path: path}
}

func (p *Picture) AddLargeDestination(path string) {
	p.Destination.Large = &PictureDestinationEntry{Path: path}
}

func (p *Picture) Prepare(variant common.PictureVariant) error {
	var (
		base64EncodedFile *string
		err               error
	)

	switch variant {
	case common.ORIGIN:
		base64EncodedFile, err = p.Destination.Origin.Tobase64()
	case common.SMALL:
		base64EncodedFile, err = p.Destination.Small.Tobase64()
	case common.MIDDLE:
		base64EncodedFile, err = p.Destination.Middle.Tobase64()
	case common.LARGE:
		base64EncodedFile, err = p.Destination.Large.Tobase64()
	}

	p.Raw = base64EncodedFile

	return err
}
