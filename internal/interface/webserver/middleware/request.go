package middleware

import (
	"context"
	"fewoserv/internal/infrastructure/common"
	"net/http"

	"github.com/google/uuid"
)

// requestMiddleware implements logic to add a request id to the http header
func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("Request-id", requestID)

		correlationID := uuid.New().String()
		ctx := context.WithValue(r.Context(), common.CorrelationIdentifier, correlationID)

		// Use the new context for the request
		newRequest := r.WithContext(ctx)

		next.ServeHTTP(w, newRequest)
	})
}
