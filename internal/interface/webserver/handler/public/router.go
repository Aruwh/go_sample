package public

import (
	applicationAdminUser "fewoserv/internal/application/admin_user"
	applicationApartment "fewoserv/internal/application/apartment"
	attributeApartment "fewoserv/internal/application/attribute"
	applicationBooking "fewoserv/internal/application/booking"
	applicationNews "fewoserv/internal/application/news"
	applicationPicture "fewoserv/internal/application/picture"
	applicationProcessLog "fewoserv/internal/application/process_log"
	applicationPublic "fewoserv/internal/application/public"
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
	processlogApp := applicationProcessLog.New(mongoDBClient)
	attributeApp := attributeApartment.New(mongoDBClient, processlogApp)
	apartmentApp := applicationApartment.New(mongoDBClient, processlogApp)
	pictureApp := applicationPicture.New(mongoDBClient, processlogApp)
	seasonApp := applicationSeason.New(mongoDBClient, processlogApp)
	settingsApp := applicationSettings.New(mongoDBClient, processlogApp)
	userApp := applicationUser.New(mongoDBClient, nil, processlogApp, nil, nil)
	adminUserApp := applicationAdminUser.New(mongoDBClient, nil, processlogApp, nil, nil)
	newsApp := applicationNews.New(mongoDBClient, processlogApp)
	bookingApp := applicationBooking.New(mongoDBClient, emailHandler, processlogApp, seasonApp, settingsApp, apartmentApp, userApp, adminUserApp, &landingpageEndpoint, copyEmailAddresses)

	publicApp := applicationPublic.New(apartmentApp, attributeApp, pictureApp, bookingApp, processlogApp, newsApp, userApp)

	handler := NewHandler(publicApp)

	// public route
	router.HandleFunc("/public/apartments", handler.GetApartments).Methods(http.MethodGet)
	router.HandleFunc("/public/apartmentDetails/{id}", handler.GetApartmentDetails).Methods(http.MethodGet)

	router.HandleFunc("/public/placeBooking", handler.PlaceBooking).Methods(http.MethodPost)
	router.HandleFunc("/public/calcBookingPrice", handler.CalculateBookingPrice).Methods(http.MethodGet)
	router.HandleFunc("/public/pictures", handler.GetPictures).Methods(http.MethodGet)
	router.HandleFunc("/public/news", handler.GetNews).Methods(http.MethodGet)
}
