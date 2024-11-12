package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/internal/infrastructure/utils"
	"net/http"
	"strings"
	"time"
)

var log = logger.New("RESPONSE")

type (
	DeletedResponse struct {
		Deleted bool `json:"deleted"`
	}

	EmailSendResponse struct {
		Send bool `json:"send"`
	}

	// Response defines the user facing response structure which will be used on each request
	Response struct {
		Data          interface{} `json:"data"`
		Error         *string     `json:"error"`
		Timestamp     int64       `json:"timestamp"`
		Domain        string      `json:"domain"`
		CorrelationID string      `json:"correlationID"`
	}
)

func extractErrorCode(err error) *string {
	delimiter := "|"

	hasErrorCode := strings.Contains(err.Error(), delimiter)
	if !hasErrorCode {
		return nil
	}

	splitedError := strings.Split(err.Error(), delimiter)
	errorCode := strings.TrimSpace(splitedError[0])

	return &errorCode
}

// CreateSuccess implements the preparation of the defined user facing success Response with the status code and the payload
func (r Response) CreateSuccess(ctx context.Context, w http.ResponseWriter, successResponseCode int, responsePayload interface{}) {
	r.Create(ctx, w, nil, &successResponseCode, responsePayload)
}

// CreateError implements the preparation of the defined user facing Response in case of failure
func (r Response) CreateError(ctx context.Context, w http.ResponseWriter, statusCode int, err error, responsePayload interface{}) {
	log.ErrorWithCtx(ctx, err.Error())

	errorCode := extractErrorCode((err))
	if errorCode != nil {
		err = errors.New(*errorCode)
	}

	r.Create(ctx, w, err, &statusCode, responsePayload)
}

// Create implements the preparation of the defined user facing Response and should be used on each response
func (r Response) Create(ctx context.Context, w http.ResponseWriter, err error, statusCode *int, responsePayload interface{}) {
	response := Response{
		Data:          nil,
		Error:         nil,
		Domain:        r.Domain,
		Timestamp:     time.Now().Unix(),
		CorrelationID: utils.ExtractCorrelationID(ctx),
	}

	if err != nil {
		strError := err.Error()
		response.Error = &strError
	}

	if responsePayload != nil {
		response.Data = responsePayload
	}

	var usedStatusCode = http.StatusOK
	if statusCode != nil {
		usedStatusCode = *statusCode
	}

	w.WriteHeader(usedStatusCode)

	// Encode the response, and check for errors.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle the encoding error, possibly by logging it or sending an error response.
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func BuildDeleteResponse(value bool) *DeletedResponse {
	deleteResponse := DeletedResponse{Deleted: value}
	return &deleteResponse
}

func BuildEmailSendResponse(value bool) *EmailSendResponse {
	emailSendResponse := EmailSendResponse{Send: value}
	return &emailSendResponse
}
