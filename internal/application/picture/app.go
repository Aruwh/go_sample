package application

import (
	"bytes"
	processLogApp "fewoserv/internal/application/process_log"
	"image"
	"image/jpeg"

	shared "fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	repository "fewoserv/internal/repository/picture"
	"fewoserv/pkg/mongodb"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/nfnt/resize"
)

var log = logger.New("APPLICATION")

type (
	IApplication interface {
		Upsert(userID, storagePath string, description *shared.Translation, file multipart.File, isOrigin *bool, recordID *string) (*shared.Picture, error)
		Delete(userID, recordID string) error
		Get(recordID string, variant common.PictureVariant) (*shared.Picture, error)
		GetMany(recordIDs []string, variant common.PictureVariant) ([]*shared.Picture, error)
		Update(userID, recordID string, variant common.PictureVariant, description *shared.Translation) (*shared.Picture, error)
	}

	Application struct {
		processLog processLogApp.IApplication
		repo       *repository.Repo
	}
)

func New(mongoDbClient mongodb.IClient, processLog processLogApp.IApplication) IApplication {
	application := Application{repo: repository.New(mongoDbClient), processLog: processLog}

	return &application
}

func ensurePathExists(storagePath string) error {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("%w: %v: %s", ErrCantCreateStorageDir, err, storagePath)
		}
	}

	return nil
}

func storeFile(file multipart.File, storagePath string) (string, *os.File, error) {
	err := ensurePathExists(storagePath)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrCantStoreFile, err)
	}
	storageDestination := fmt.Sprintf("%s/%d%d.jpg", storagePath, time.Now().Unix(), &file)

	newFile, err := os.Create(storageDestination)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v: %s", ErrCantStoreFile, err, storageDestination)
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v: %s", ErrCantStoreFile, err, storageDestination)

	}

	return storageDestination, newFile, nil
}

func storeImage(img image.Image, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("%w: %v: %s", ErrCantDeleteFile, err, filePath)
	}

	return nil
}

func deleteFiles(destination shared.PictureDestination) error {
	var err error

	err = deleteFile(destination.Origin.Path)
	if err != nil {
		return err
	}
	err = deleteFile(destination.Large.Path)
	if err != nil {
		return err
	}
	err = deleteFile(destination.Middle.Path)
	if err != nil {
		return err
	}
	err = deleteFile(destination.Small.Path)
	if err != nil {
		return err
	}

	return nil
}

func resizeImage(width, height uint, imageFile *os.File) (image.Image, error) {
	fileBytes, err := io.ReadAll(imageFile)
	if err != nil {
		return nil, err
	}

	image, err := jpeg.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	resizedImage := resize.Resize(width, height, image, resize.NearestNeighbor)

	return resizedImage, nil
}

func resizeImageProportional(maxWidth, maxHeight uint, imageFile *os.File) (image.Image, error) {
	_, err := imageFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	originalWidth := uint(img.Bounds().Dx())
	originalHeight := uint(img.Bounds().Dy())

	var newWidth, newHeight uint
	if originalWidth > originalHeight {
		newWidth = maxWidth
		scaleFactor := float64(maxWidth) / float64(originalWidth)
		newHeight = uint(float64(originalHeight) * scaleFactor)
	} else {
		newHeight = maxHeight
		scaleFactor := float64(maxHeight) / float64(originalHeight)
		newWidth = uint(float64(originalWidth) * scaleFactor)
	}

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.NearestNeighbor)

	return resizedImage, nil
}

func resizeAndStoreImage(width, height int, file *os.File, storagePath string) (string, error) {
	resizedImage, err := resizeImageProportional(uint(width), uint(height), file)
	if err != nil {
		return "", err
	}

	err = ensurePathExists(storagePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCantStoreFile, err)
	}
	storageDestination := fmt.Sprintf("%s/%d%d.jpg", storagePath, time.Now().Unix(), &file)

	err = storeImage(resizedImage, storageDestination)
	if err != nil {
		return "", err
	}

	return storageDestination, nil
}

func (a *Application) Upsert(userID, storagePath string, description *shared.Translation, file multipart.File, isOrigin *bool, recordID *string) (*shared.Picture, error) {
	usedPicture := shared.NewPicture(description)

	// we overwrite the pciture with an already existing one.
	// it has the background that we use two different image formats, which are realised in two upload processes
	if recordID != nil {
		loadedPicture, err := a.repo.LoadByID(*recordID)
		if err == nil {
			usedPicture = loadedPicture
		}
	}

	////////////////////////////////////////
	// original
	////////////////////////////////////////
	originImgPath := fmt.Sprintf("%s/origin", storagePath)
	originalFileStoragePath, originalFile, err := storeFile(file, originImgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}
	usedPicture.AddOriginDestination(originalFileStoragePath)
	a.processLog.New(userID, originalFileStoragePath, common.CREATED, common.PICTURE, &usedPicture.ID)

	////////////////////////////////////////
	// large
	////////////////////////////////////////
	if isOrigin != nil && *isOrigin == true {
		largeImgPath := fmt.Sprintf("%s/large", storagePath)
		largeImgStoragePath, err := resizeAndStoreImage(1600, 1200, originalFile, largeImgPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
		}
		usedPicture.AddLargeDestination(largeImgStoragePath)
		a.processLog.New(userID, largeImgStoragePath, common.CREATED, common.PICTURE, &usedPicture.ID)

		usedPicture.Prepare(common.MIDDLE)
		err = a.repo.Upsert(usedPicture)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
		}

		return usedPicture, nil
	}

	////////////////////////////////////////
	// middle
	////////////////////////////////////////
	middleImgPath := fmt.Sprintf("%s/middle", storagePath)
	middleImgStoragePath, err := resizeAndStoreImage(800, 400, originalFile, middleImgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}
	usedPicture.AddMiddleDestination(middleImgStoragePath)
	a.processLog.New(userID, middleImgStoragePath, common.CREATED, common.PICTURE, &usedPicture.ID)

	////////////////////////////////////////
	// small
	////////////////////////////////////////
	smallImgPath := fmt.Sprintf("%s/small", storagePath)
	smallImgStoragePath, err := resizeAndStoreImage(400, 200, originalFile, smallImgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}
	usedPicture.AddSmallnDestination(smallImgStoragePath)
	a.processLog.New(userID, smallImgStoragePath, common.CREATED, common.PICTURE, &usedPicture.ID)

	err = a.repo.Upsert(usedPicture)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantSave, err)
	}

	usedPicture.Prepare(common.MIDDLE)

	return usedPicture, nil
}

func (a *Application) Delete(userID, recordID string) error {
	pictureToDelete, err := a.repo.LoadByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = a.repo.DeleteByID(recordID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	err = deleteFiles(*pictureToDelete.Destination)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantDelete, err)
	}

	a.processLog.New(userID, pictureToDelete.Destination.Origin.Path, common.DELETED, common.PICTURE, &recordID)
	a.processLog.New(userID, pictureToDelete.Destination.Large.Path, common.DELETED, common.PICTURE, &recordID)
	a.processLog.New(userID, pictureToDelete.Destination.Middle.Path, common.DELETED, common.PICTURE, &recordID)
	a.processLog.New(userID, pictureToDelete.Destination.Small.Path, common.DELETED, common.PICTURE, &recordID)

	return nil
}

func (a *Application) Get(recordID string, variant common.PictureVariant) (*shared.Picture, error) {
	foundPicture, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	err = foundPicture.Prepare(variant)
	if err != nil {
		return nil, fmt.Errorf("%w: %v: %v: %s", ErrRecordNotExists, ErrCantTransformToBase64, err, recordID)
	}

	return foundPicture, nil
}

func (a *Application) GetMany(recordIDs []string, variant common.PictureVariant) ([]*shared.Picture, error) {
	foundPictures, err := a.repo.FindMany(recordIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	for _, picture := range foundPictures {
		err = picture.Prepare(variant)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: %v: %s", ErrRecordNotExists, ErrCantTransformToBase64, err, picture.ID)
		}
	}

	return foundPictures, nil
}

func (a *Application) Update(userID, recordID string, variant common.PictureVariant, description *shared.Translation) (*shared.Picture, error) {
	foundPicture, err := a.repo.LoadByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecordNotExists, err)
	}

	foundPicture.Update(description)

	err = a.repo.Update(foundPicture)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUpdate, err)
	}

	err = foundPicture.Prepare(variant)
	if err != nil {
		return nil, fmt.Errorf("%w: %v: %v: %s", ErrCantUpdate, ErrCantTransformToBase64, err, recordID)
	}

	a.processLog.New(userID, "", common.UPDATED, common.PICTURE, &recordID)

	return foundPicture, nil
}
