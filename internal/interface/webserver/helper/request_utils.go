package helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	// SessionTokenKey - key containing UAT
	SessionTokenKey = "SESSION_TOKEN"
)

// GetStringFromFormField implements logic to extract a parameter from the form field
func GetStringFromFormField(r *http.Request, fieldName string) *string {
	values := r.MultipartForm.Value[fieldName]
	if len(values) > 0 {
		return &values[0]
	}
	return nil
}

// GetQueryParam implements logic to extract a parameter from the request URL
func GetQueryParam(r *http.Request, fieldName string) *string {
	value := r.URL.Query().Get(fieldName)
	if len(value) == 0 {
		fmt.Printf("requested value of fieldName: %s not available\n", fieldName)
		return nil
	}

	return &value
}

// GetQueryVar implements logic to extract a variable from the request URL
func GetQueryVar(r *http.Request, fieldName string) *string {
	varParams := mux.Vars(r)

	value, ok := varParams[fieldName]
	if !ok {
		fmt.Printf("requested value of fieldName: %s not available\n", fieldName)
		return nil
	}

	return &value
}

// ExtractToken func. realised to extract the passt bearer token from the request header
func ExtractToken(r *http.Request) string {
	return r.Header.Get("authorization")
}

// ExtractTokenFromCookie - extracts token from cookie
func ExtractTokenFromCookie(r *http.Request) (string, error) {
	var (
		cookie *http.Cookie
		err    error
	)

	if cookie, err = r.Cookie(SessionTokenKey); err != nil {
		return "", err
	}

	return strings.Replace(cookie.Value, "Bearer ", "", -1), nil
}
