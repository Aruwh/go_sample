package booking

import (
	applicationAdminUser "fewoserv/internal/application/admin_user"
	applicationApartment "fewoserv/internal/application/apartment"
	application "fewoserv/internal/application/booking"
	applicationProcessLog "fewoserv/internal/application/process_log"
	applicationSeason "fewoserv/internal/application/saison"
	applicationSettings "fewoserv/internal/application/settings"
	applicationUser "fewoserv/internal/application/user"
	emailhandler "fewoserv/internal/interface/email_handler"

	"fewoserv/pkg/mongodb"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRouter registered all real_estate endpoints
func RegisterRouter(router *mux.Router, mongoDBClient mongodb.IClient, emailHandler emailhandler.IEmailHandler, landingpageEndpoint string, copyEmailAddresses []string) {
	appProcessLog := applicationProcessLog.New(mongoDBClient)

	appSeason := applicationSeason.New(mongoDBClient, appProcessLog)
	appSettings := applicationSettings.New(mongoDBClient, appProcessLog)
	appApartment := applicationApartment.New(mongoDBClient, appProcessLog)
	appUser := applicationUser.New(mongoDBClient, nil, appProcessLog, nil, nil)
	appAdminUser := applicationAdminUser.New(mongoDBClient, nil, appProcessLog, nil, nil)

	app := application.New(mongoDBClient, emailHandler, appProcessLog, appSeason, appSettings, appApartment, appUser, appAdminUser, &landingpageEndpoint, copyEmailAddresses)
	handler := NewHandler(app)

	router.HandleFunc("/bookings", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/bookings/{id}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/bookings/{id}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/bookings/{id}", handler.Update).Methods(http.MethodPatch)
	router.HandleFunc("/bookings/{id}/message", handler.AddMessage).Methods(http.MethodPatch)
	router.HandleFunc("/bookings", handler.GetBookingOverviews).Methods(http.MethodGet)
}
