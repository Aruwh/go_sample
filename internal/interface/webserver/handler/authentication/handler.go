package authentication

import (
	application "fewoserv/internal/application/authentication"
	"fewoserv/internal/infrastructure/logger"
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/internal/interface/webserver/helper"
	"fewoserv/internal/interface/webserver/shared"
	"fewoserv/pkg/mongodb"
	"net/http"
)

var (
	DOMAIN = "Authentication"
	log    = logger.New("AUTH")
)

type (
	// IHandler Interface defines the structure of functions which needed to cover all route implementations
	IHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
		ForgotPwd(w http.ResponseWriter, r *http.Request)
		ResetPwd(w http.ResponseWriter, r *http.Request)
	}

	// Handler defines the used handler which combines all implementations of available endpoints
	Handler struct {
		response                          shared.Response
		mongoDBClient                     mongodb.IClient
		emailHandler                      emailhandler.IEmailHandler
		feEndpoint                        string
		jwtExpireTimeInMinutes            int
		jwtExpireTimeForPwdResetInMinutes int
	}
)

// NewHandler creates a new Handler
func NewHandler(jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes int, feEndpoint string, mongoDbClient mongodb.IClient, emailHandler emailhandler.IEmailHandler) IHandler {
	var handler = Handler{
		jwtExpireTimeInMinutes:            jwtExpireTimeInMinutes,
		jwtExpireTimeForPwdResetInMinutes: jwtExpireTimeForPwdResetInMinutes,
		feEndpoint:                        feEndpoint,
		emailHandler:                      emailHandler,
		mongoDBClient:                     mongoDbClient}

	handler.response.Domain = DOMAIN

	return handler
}

// Get implements logic to get health information of the mono api and the searcher service
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var requestData LoginAdminRequest
	err := helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var token string
	token, err = application.Login(h.mongoDBClient, requestData.Email, requestData.RawPassword, h.jwtExpireTimeInMinutes)

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusUnauthorized, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, LoginAdminResponse{}.ConvertFromValue(token))
}

// Get implements logic to get health information of the mono api and the searcher service
func (h Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) ForgotPwd(w http.ResponseWriter, r *http.Request) {
	var requestData ForgotPwdRequest
	err := helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	err = application.ForgotPwd(h.mongoDBClient, h.emailHandler, requestData.Email, h.feEndpoint, h.jwtExpireTimeForPwdResetInMinutes)

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusForbidden, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, ForgotPwdResponse{}.ConvertFromValue(true))
}

func (h Handler) ResetPwd(w http.ResponseWriter, r *http.Request) {
	var requestData ResetPwdRequest
	err := helper.AssignAndValidateJSON(&requestData, r.Body)
	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusBadRequest, err, nil)
		return
	}

	var token string
	token, err = application.ResetPwd(h.mongoDBClient, requestData.Pwd, requestData.RepeatedPwd, requestData.Token, h.jwtExpireTimeInMinutes)

	if err != nil {
		h.response.CreateError(r.Context(), w, http.StatusUnauthorized, err, nil)
		return
	}

	h.response.CreateSuccess(r.Context(), w, http.StatusOK, LoginAdminResponse{}.ConvertFromValue(token))
}
