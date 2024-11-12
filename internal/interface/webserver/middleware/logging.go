package middleware

import (
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/interface/webserver/helper"
	"fmt"
	"net/http"
)

var log = logger.New("WEBSERVER")

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callingIP := helper.GetIP(r)

		log.InfoWithCtx(r.Context(), fmt.Sprintf("%s %s	%s", r.Method, r.RequestURI, callingIP))

		next.ServeHTTP(w, r)
	})
}
