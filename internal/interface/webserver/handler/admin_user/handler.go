package adminuser

import (
	application "fewoserv/internal/application/admin_user"
	adminuser "fewoserv/internal/domain/admin_user"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"fmt"
	"net/http"
	"strconv"
)

var (
	DOMAIN = "AdminUser"
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		GetMe(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
		GetMany(w http.ResponseWriter, r *http.Request)
		Create(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		Invite(w http.ResponseWriter, r *http.Request)
	}

	// HealthHandler defines the used handler which combines all implementations of available endpoints
	HealthHandler struct {
		response    shared.Response
		application application.IApplication
	}
)

// NewHandler creates a new HealthHandler
func NewHandler(applicationAdminUser application.IApplication) IHandler {
	var handler = HealthHandler{}

	handler.response.Domain = DOMAIN
	handler.application = applicationAdminUser

	return handler
}

func (h HealthHandler) Create(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData CreateRequest
	err = helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	randomPassword := utils.GenerateRandomString(24)

	var adminUser *adminuser.AdminUser
	switch requestData.Type {
	case common.ADMINISTRATOR:
		adminUser, err = h.application.CreateAdmin(identity.UserID, requestData.Email, randomPassword, randomPassword, requestData.FirstName, requestData.LastName, requestData.Locale)
	case common.APARTMENT_OWNER:
		adminUser, err = h.application.CreateApartmentOwner(identity.UserID, requestData.Email, randomPassword, randomPassword, requestData.FirstName, requestData.LastName, requestData.Locale)
	default:
		err = fmt.Errorf("%w: %v", shared.ErrUnsupportedAdminType, requestData.Type)
	}

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, adminUser)
}

func (h HealthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission()
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	adminUser, err := h.application.Get(identity.UserID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, adminUser)
}

func (h HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	adminUserID := helper.GetQueryVar(r, "id")

	adminUser, err := h.application.Get(*adminUserID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, adminUser)
}

func (h HealthHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_VIEW)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var requestData common.GetManyRequest[GetManyFilter]

	adminUserType := helper.GetQueryParam(r, "type")
	if adminUserType != nil {
		transformedType, err := strconv.Atoi(*adminUserType)
		if err != nil {
			h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
			return
		}
		usableAdminUserType := common.AdminUserType(transformedType)
		requestData.Filter.Type = &usableAdminUserType
	}

	name := helper.GetQueryParam(r, "name")
	if name != nil {
		requestData.Filter.Name = name
	}

	field := helper.GetQueryParam(r, "field")
	requestData.Sort.Field = common.SortByType(*field)

	order := helper.GetQueryParam(r, "order")
	requestData.Sort.Order = common.OrderType(*order)

	limit, _ := strconv.ParseInt(*helper.GetQueryParam(r, "limit"), 10, 64)
	requestData.Limit = limit

	skip, _ := strconv.ParseInt(*helper.GetQueryParam(r, "skip"), 10, 64)
	requestData.Skip = skip

	adminUsers, err := h.application.GetMany(requestData.Filter.Type, requestData.Filter.Name, requestData.Sort, requestData.Skip, requestData.Limit)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, adminUsers)
}

func (h HealthHandler) Update(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_EDIT)
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

	updateAdminUserID := helper.GetQueryVar(r, "id")
	transformedPasswordUpdate := TransformRequestUpdatePassword(requestData.PasswordUpdate)

	adminUser, err := h.application.Update(identity.UserID, *updateAdminUserID, requestData.FirstName, requestData.LastName, requestData.IsActive, requestData.Type, requestData.Permissions, transformedPasswordUpdate, requestData.Locale)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, adminUser)
}

func (h HealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_DELETE)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	deleteAdminUserID := helper.GetQueryVar(r, "id")

	err = h.application.Delete(identity.UserID, *deleteAdminUserID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, shared.BuildDeleteResponse(false))
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, shared.BuildDeleteResponse(true))
}

func (h HealthHandler) Invite(w http.ResponseWriter, r *http.Request) {
	identity := helper.ExtractIdentity(r)
	err := identity.EnsureRequestPermission(common.ADMIN_USER_EDIT)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	adminUserID := helper.GetQueryVar(r, "id")

	err = h.application.SendInvitation(identity.UserID, *adminUserID)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, shared.BuildEmailSendResponse(true))
}
