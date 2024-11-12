package middleware

import (
	"encoding/json"
	"errors"
	"fewoserv/internal/infrastructure/utils"
	"fewoserv/internal/interface/webserver/shared"
	"fmt"
	"net/http"
	"time"
)

// recoveryMiddleware implements logic to catch a unhandled panic/exception and to provide a proper response to the client in such cases
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			log.ErrorWithCtx(r.Context(), fmt.Sprintf("An unhandled error occurred REQUEST: [%+v](%+v) ERROR: %+v", r.Method, r.URL.String(), err))

			// To convert the "err" interface could send some confusing messages to the client, even if it was of type "error"
			returnedError := errors.New("an unhandled error occurred").Error()
			responseBody, _ := json.Marshal(shared.Response{
				Data: nil,
				// We cant't reliably get the requested domain from here, so we use unknown.
				// Error:     processerrors.NewUnknownError(returnedError, "UNKNOWN"),
				Error:         &returnedError,
				Timestamp:     time.Now().Unix(),
				CorrelationID: utils.ExtractCorrelationID(r.Context()),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(responseBody)
		}()

		next.ServeHTTP(w, r)
	})
}
