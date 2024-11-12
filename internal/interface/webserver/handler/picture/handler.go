package picture

import (
	"encoding/json"
	application "fewoserv/internal/application/picture"
	sharedDomain "fewoserv/internal/domain/shared"
	"fewoserv/internal/infrastructure/common"
	picturecache "fewoserv/internal/infrastructure/picture_cache"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"fmt"
	"net/http"
)

var (
	DOMAIN = "Picture"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Upload(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetMany(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
	}

	// HealthHandler defines the used handler which combines all implementations of available endpoints
	HealthHandler struct {
		response     shared.Response
		application  application.IApplication
		storagePath  string
		pictureCache picturecache.IPictureCache
	}
)

// NewHandler creates a new HealthHandler
func NewHandler(applicationPicture application.IApplication, storagePath string, pictureCache picturecache.IPictureCache) IHandler {
	var handler = HealthHandler{}

	handler.response.Domain = DOMAIN
	handler.application = applicationPicture
	handler.storagePath = storagePath
	handler.pictureCache = pictureCache

	return handler
}

func (h HealthHandler) Upload(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PICTURE_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	// Parse the multipart request
	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		description sharedDomain.Translation
		isOrigin    *bool
		recordID    *string
	)

	extractedDescription := helper.GetStringFromFormField(r, "description")
	// Unmarshal the JSON string into the Translation struct
	err = json.Unmarshal([]byte(*extractedDescription), &description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extractedOrigin := helper.GetStringFromFormField(r, "isOrigin")
	// Unmarshal the JSON string into the Translation struct
	err = json.Unmarshal([]byte(*extractedOrigin), &isOrigin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recordID = helper.GetStringFromFormField(r, "recordID")

	createdPictures := []*sharedDomain.Picture{}
	// Iterate through all files in the request
	for multipartFileIdentifier, _ := range r.MultipartForm.File {
		multipartFile, _, err := r.FormFile(multipartFileIdentifier)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer multipartFile.Close()

		createdPicture, err := h.application.Upsert(identity.UserID, h.storagePath, &description, multipartFile, isOrigin, recordID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// the upload process of an image can be done in two steps. the reason for this is that two image formats were desired, but there was no complete revision of the upload due to time constraints
		if isOrigin != nil && *isOrigin == true {
			h.pictureCache.Add("large", createdPicture)
		}

		createdPictures = append(createdPictures, createdPicture)
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, createdPictures)
}

func (h HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PICTURE_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	recordID := helper.GetQueryVar(r, "id")
	variant := helper.GetQueryParam(r, "variant")

	cachedPicture := h.pictureCache.Get(*variant, *recordID)
	if cachedPicture != nil {
		h.response.CreateSuccess(r.Context(), w, http.StatusOK, cachedPicture)
		return
	}

	picture, err := h.application.Get(*recordID, common.PictureVariant(*variant))
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.pictureCache.Add(*variant, picture)

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, picture)

}

func (h HealthHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PICTURE_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}
}

func (h HealthHandler) Update(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PICTURE_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData UpdateRequest
	err = helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	updatePictureID := helper.GetQueryVar(r, "id")
	variant := helper.GetQueryParam(r, "variant")

	transformedDescription := shared.TransformRequestTranslation(requestData.Description)

	picture, err := h.application.Update(identity.UserID, *updatePictureID, requestData.Variant, transformedDescription)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.pictureCache.Add(*variant, picture)

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, picture)
}

func (h HealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.PICTURE_DELETE)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	deletePictureID := helper.GetQueryVar(r, "id")

	err = h.application.Delete(identity.UserID, *deletePictureID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, shared.BuildDeleteResponse(true))
}
