package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"fewoserv/internal/infrastructure/common"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func SymmetricDifference[T comparable](slice1, slice2 []T) []T {
	elements := make(map[T]int)

	for _, elem := range slice1 {
		elements[elem]++
	}

	for _, elem := range slice2 {
		elements[elem]++
	}

	var symmetricDiff []T

	for elem, count := range elements {
		if count == 1 {
			symmetricDiff = append(symmetricDiff, elem)
		}
	}

	return symmetricDiff
}

func Intersection[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	result := make([]T, 0)

	// Add elements of slice2 to set
	for _, value := range slice2 {
		set[value] = true
	}

	// Check for difference in slice1
	for _, value := range slice1 {
		if !set[value] {
			result = append(result, value)
		}
	}

	return result
}

func DecryptPwd(password string) (*string, error) {
	hashedPwdBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hashedPwdBytes)
	return &hashedPassword, nil
}

func ValidatePwdValid(password, compareRawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(compareRawPassword))
}

func ValidateEmail(email string) bool {
	// Define a regular expression for a basic email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Use the MatchString method to check if the email matches the pattern
	return re.MatchString(email)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Create a byte slice to hold the random bytes
	randomBytes := make([]byte, length)

	// Use crypto/rand to generate random bytes
	for i := 0; i < length; i++ {
		randomByte, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Println("Error generating random string:", err)
			return ""
		}
		randomBytes[i] = charset[randomByte.Int64()]
	}

	// Convert the random bytes to a string
	randomString := string(randomBytes)

	return randomString
}

func isValidPassword(password string) bool {
	// Check if the password is at least 8 characters long
	if len(password) < 8 {
		return false
	}

	// Check if the password contains at least one uppercase letter, one lowercase letter, one digit, and one special character
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsLower(char) {
			hasLowercase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialChar = true
		}
	}

	if !(hasUppercase && hasLowercase && hasDigit && hasSpecialChar) {
		return false
	}

	return true
}

func SwapDatesIfNeeded(fromDate, toDate time.Time) (time.Time, time.Time) {
	swapedFromDate := fromDate
	swapedToDate := toDate

	shouldBeSwaped := fromDate.After(toDate) || toDate.Before(fromDate)
	if shouldBeSwaped {
		swapedFromDate = toDate
		swapedToDate = fromDate
	}

	return swapedFromDate, swapedToDate
}

func ExtractCorrelationID(ctx context.Context) string {
	value := ctx.Value(common.CorrelationIdentifier)
	return value.(string)
}

func ConvertUnixStringToInt(strPtr *string) (int, error) {
	if strPtr == nil {
		return 0, errors.New("passed value is nil")
	}

	intValue, err := strconv.Atoi(*strPtr)
	if err != nil {
		return 0, err
	}

	// transform unix date with milliseconds
	if len(*strPtr) == 13 {
		intValue = intValue / 1000
	}

	return intValue, nil
}
