package middleware

import (
	"context"
	"errors"
	"fewoserv/internal/infrastructure/common"
	helper "fewoserv/internal/interface/webserver/helper/authentication"
	"fewoserv/internal/interface/webserver/shared"
	"net/http"
	"strings"
)

func extractToken(r *http.Request) (string, bool) {
	hasToken := false
	bearerToken := r.Header.Get("Authorization")

	var tokenString = ""
	if len(bearerToken) != 0 {
		tokenString = bearerToken[7:]
		hasToken = len(tokenString) > 0
	}

	return tokenString, hasToken
}

func getUri(r *http.Request) string {
	leadingURIPart := strings.Split(r.RequestURI, "?")

	return leadingURIPart[0]
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, hasToken := extractToken(r)

		if hasToken {
			token, err := helper.ValidateAndTransformToken(tokenString)
			if err != nil {
				shared.Response{}.CreateError(r.Context(), w, http.StatusUnauthorized, err, nil)
				return
			}

			identity := helper.TransformTokenToIdentity(token)

			// Create a context and add identity
			ctx := context.WithValue(r.Context(), common.IdentityIdentifier, *identity)

			// Use the new context for the request
			newRequest := r.WithContext(ctx)

			next.ServeHTTP(w, newRequest)
			return
		}

		// exception for open routes which don't need an authentication
		uri := getUri(r)
		if strings.Contains(uri, "/public/apartmentDetails") {
			next.ServeHTTP(w, r)
			return
		}

		switch uri {
		case "/login", "/forgotPwd", "/resetPwd", "/ping", "/healthz":
			next.ServeHTTP(w, r)
			return

			// public routes for landing page
		case "/public/apartments", "/public/placeBooking", "/public/calcBookingPrice", "/public/pictures", "/public/news":
			next.ServeHTTP(w, r)
			return

		}

		shared.Response{}.CreateError(r.Context(), w, http.StatusUnauthorized, errors.New("shit happens ... no token send"), nil)
	})
}
