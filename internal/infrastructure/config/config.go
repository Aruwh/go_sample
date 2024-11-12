package config

import (
	"errors"
	"fewoserv/internal/infrastructure/common"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

var (
	cfg *Config
)

// Config contains every single configuration variable for the application to run.
// These variables should either be set in .env file or as system variables.
type (
	Config struct {
		Service struct {
			LogLevel               common.LogLevel `envconfig:"LOG_LEVEL"`
			HTTPPort               string          `envconfig:"HTTP_PORT"`
			JwtSecret              string          `envconfig:"JWT_SECRET"`
			JwtExpireTimeInMinutes int             `envconfig:"JWT_EXPIRE_TIME_IN_MINUTES"`
			CORS_DEBUGING          bool            `envconfig:"CORS_DEBUGING"`
			CORSAllowedOrigins     []string        `envconfig:"CORS_ALLOWED_ORIGINS"`
			StoragePath            string          `envconfig:"STORAGE_PATH"`
		}

		MongoDB struct {
			MongoDBUri  string `envconfig:"MONGODB_URI"`
			MongoDBName string `envconfig:"MONGODB_DB_NAME"`
		}

		Email struct {
			From                string   `envconfig:"MAIL_FROM"`
			Password            string   `envconfig:"MAIL_PASSWORD"`
			ServerURL           string   `envconfig:"MAIL_SERVER_URL"`
			ServerPort          string   `envconfig:"MAIL_SERVER_PORT"`
			FeEndpoint          string   `envconfig:"FE_ENDPOINT"`
			LandingpageEndpoint string   `envconfig:"LANDINGPAGE_ENDPOINT"`
			CopyEmailAddresses  []string `envconfig:"MAIL_COPY_ADDRESSES"`
		}

		Authentication struct {
			JwtExpireTimeForPwdResetInMinutes int `envconfig:"JWT_EXPIRE_TIME_FOR_PWD_RESET_IN_MINUTES"`
		}

		Ssl struct {
			Enabled  string `envconfig:"SSL_ENABLED,optional"`
			CertPath string `envconfig:"SSL_CERT_PATH,optional"`
			KeyPath  string `envconfig:"SSL_KEY_PATH,optional"`
		}
	}
)

// findEnv is used to find the .env by running up folders until one is found.
func findEnv() string {
	currentWorkDirectory, _ := os.Getwd()
	foundPath := ""
	path := `.env`
	// travel over each folder level and try to find the .env
	for level := 1; level < 20; level++ {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			// prepare next path
			path = currentWorkDirectory + "/" + strings.Repeat("../", level) + `.env`
			// path does not exist
			continue
		}
		foundPath = path
		break
	}
	return foundPath
}

// Load is returning values from the .env configuration.
func Load() *Config {
	if cfg != nil {
		return cfg
	}

	path := findEnv()
	_ = godotenv.Load(path) // We ignore this error. This method throws an error if the .env file could not be found.

	if err := envconfig.Init(&cfg); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	return cfg
}
