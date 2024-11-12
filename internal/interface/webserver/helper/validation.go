package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

/*
AssignAndValidateJSON is used to decode a binary array with encoding/json into a struct and validate the struct with github.com/go-playground/validator.
NOTE: The struct must be passed as a reference.
*/
func AssignAndValidateJSON[T any](target *T, reader io.Reader) error {
	jsonBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("JSON failed to ready payload")
		return errors.New("failed to ready payload")
	}

	jsonReader := bytes.NewReader(jsonBytes)
	tokenDecoder := json.NewDecoder(jsonReader)

	// Check if the provided value was null, null is technically a json but not allowed here
	// this also returns an error if the provided data is garbage (not a valid json)
	jsonToken, err := tokenDecoder.Token()
	if jsonToken == nil || err != nil {
		fmt.Printf("JSON invalid payload: %s", jsonToken)
		return errors.New("invalid payload")
	}

	// Reset the reader and read the entire json, now that we know its not null
	_, seekError := jsonReader.Seek(0, io.SeekStart)
	if seekError != nil {
		fmt.Printf("JSON failed to ready payload")
		return errors.New("failed to ready payload")
	}

	decoder := json.NewDecoder(jsonReader)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(target)
	if err != nil {
		fmt.Printf("JSON decode Error: %+v\n", err)
		return err
	}

	validate := validator.New()
	validate.RegisterValidation("mongoDbID", validateMongoDBID)
	if err := validate.Struct(target); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, fieldErr := range errs {
			fmt.Printf("field %s: %s\n", fieldErr.Field(), fieldErr.Tag())
		}
		return err
	}

	return nil
}

func validateMongoDBID(fl validator.FieldLevel) bool {
	// Regular expression for a valid MongoDB ObjectId
	pattern := `^[0-9a-fA-F]{24}$`

	return regexp.MustCompile(pattern).MatchString(fl.Field().String())
}
