package helper

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/config"
	"fewoserv/internal/infrastructure/logger"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	log    = logger.New("WEBSERVER")
	secret = []byte(config.Load().Service.JwtSecret)
)

func getAudience(adminUserType common.AdminUserType) string {
	var audience = common.AUDIANCE_USER

	switch adminUserType {
	case common.ADMINISTRATOR:
		audience = common.AUDIANCE_ADMIN_USER

	case common.APARTMENT_OWNER:
		audience = common.AUDIANCE_APARTMENT_OWNER

	case common.SUPER_ADMINISTRATOR:
		audience = common.AUDIANCE_SUPER_ADMIN_USER

	}

	return string(audience)
}

func GenerateJwt(userID *string, adminUserType common.AdminUserType, permissions []common.RequestPermission, locale common.Locale, expirationTimeInMinutes *int) (string, error) {
	// default 1 hour
	var usedExpirationTimeInSeconds = time.Minute * time.Duration(60)
	if expirationTimeInMinutes != nil {
		usedExpirationTimeInSeconds = time.Minute * time.Duration(*expirationTimeInMinutes)
	}

	claims := jwt.MapClaims{
		// custom stuff: permissions of the current user
		"permissions": permissions,
		// custom stuff: locale of the current user
		"locale": locale,
		// (Subject): The identity of the user or subject of the token.
		"sub": *userID,
		// (Audience): The intended audience of the token, i.e., for whom the token is meant.
		"aud": getAudience(adminUserType),
		// (Issuer): The entity that issued the token, typically the authentication authority.
		"iss": "fewoserv.com",
		// (Issued At): The time when the token was issued.
		"iat": time.Now().Unix(),
		// (Expiration): The expiration time of the token as a Unix timestamp.
		"exp": time.Now().Add(usedExpirationTimeInSeconds).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJwtForPwdReset(userID *string, adminUserType common.AdminUserType, locale common.Locale, expirationTimeInHours *int) (string, error) {
	// default 48 hour
	var usedExpirationTimeInSeconds = time.Minute * time.Duration(2880)
	if expirationTimeInHours != nil {
		usedExpirationTimeInSeconds = time.Minute * time.Duration(*expirationTimeInHours)
	}

	claims := jwt.MapClaims{
		// custom stuff: locale of the current user
		"locale": locale,
		// (Subject): The identity of the user or subject of the token.
		"sub": *userID,
		// (Audience): The intended audience of the token, i.e., for whom the token is meant.
		"aud": getAudience(adminUserType),
		// (Issuer): The entity that issued the token, typically the authentication authority.
		"iss": "fewoserv.com",
		// (Issued At): The time when the token was issued.
		"iat": time.Now().Unix(),
		// (Expiration): The expiration time of the token as a Unix timestamp.
		"exp": time.Now().Add(usedExpirationTimeInSeconds).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAndTransformToken(tokenString string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if err != nil {
		log.Error(err.Error())

		if !token.Valid {
			return nil, ErrNotAuthorised
		}

		return nil, ErrMalformedToken
	}

	//ErrTokenInvalidClaims
	if !token.Valid {
		return nil, ErrNotAuthorised
	}

	return token, nil
}
