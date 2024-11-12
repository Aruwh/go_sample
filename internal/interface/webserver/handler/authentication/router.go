package authentication

import (
	emailhandler "fewoserv/internal/interface/email_handler"
	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all healthz endpoints
func RegisterRouter(router *mux.Router, jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes int, feEndpoint string, mongoDBClient mongodb.IClient, emailHandler emailhandler.IEmailHandler) {
	handler := NewHandler(jwtExpireTimeInMinutes, jwtExpireTimeForPwdResetInMinutes, feEndpoint, mongoDBClient, emailHandler)

	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/forgotPwd", handler.ForgotPwd).Methods(http.MethodPost)
	router.HandleFunc("/resetPwd", handler.ResetPwd).Methods(http.MethodPost)
}
