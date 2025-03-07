// The config package allows an easy way to provide variable options to the
// program read at runtime during startup. These are typically settings like the
// http port, database connection string, or password.
package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	SESSION_COOKIE_NAME            = "_maxwarden_session"
	SESSION_COOKIE_EXPIRY_DAYS int = 100
	SESSION_COOKIE_ENTROPY     int = 33

	IDENTITY_COOKIE_NAME        string = "_maxwarden_identity"
	IDENTITY_COOKIE_EXPIRY_DAYS int    = 30
	IDENTITY_TOKEN_EXPIRY_DAYS  int    = 30
	IDENTITY_COOKIE_ENTROPY     int    = 33
	IDENTITY_LOGIN_PATH         string = "/auth/login"
	IDENTITY_LOGOUT_PATH        string = "/auth/logout"
	IDENTITY_DEFAULT_PATH       string = "/app"
	IDENTITY_AUTH_REDIRECT      bool   = true
	IDENTITY_AUTH_KEY           string = "CORRECT_HORSE_BATTERY_STAPLE"

	// This key is NOT used for the hashing of passwords, or secure session data over the wire.
	// It is ONLY used for performing quick file and string hashes, where security is not a factor.
	DATA_HASH_KEY string = "01234567890123456789012345678901"

	PASSWORD_MIN_LENGTH         int = 8
	PASSWORD_REQUIRED_UPPERCASE int = 1
	PASSWORD_REQUIRED_LOWERCASE int = 1
	PASSWORD_REQUIRED_NUMBERS   int = 1
	PASSWORD_REQUIRED_SYMBOLS   int = 0

	MAX_LOGIN_ATTEMPTS int = 5
)

type configuration struct {
	Domain                  string `env:"DOMAIN"`
	Host                    string `env:"HOST"`
	Port                    string `env:"PORT"`
	IdentityPrivateKey      string `env:"IDENTITY_PRIVATE_KEY"`
	IdentityDefaultPassword string `env:"IDENTITY_DEFAULT_PASSWORD"`
	SessionPrivateKey       string `env:"SESSION_PRIVATE_KEY"`
	SmtpServer              string `env:"SMTP_SERVER"`
	SmtpPort                string `env:"SMTP_PORT"`
	SmtpUsername            string `env:"SMTP_USERNAME"`
	SmtpDisplayFrom         string `env:"SMTP_DISPLAY_FROM"`
	SmtpPassword            string `env:"SMTP_PASSWORD"`
	SmtpRequireAuth         bool   `env:"SMTP_REQUIRE_AUTH"`
}

//mysql:    "DbConnectionString": "root:PASSWORD@tcp(localhost:3306)/example?parseTime=true",

var config configuration

func GetConfig() configuration {
	return config
}

func Init() {
	// When in debug mode, set environment variables from the `.env` file directly.
	// Just a developer convenience.
	if DEBUG {
		godotenv.Load()
	}

	envErr := env.Parse(&config)
	if envErr != nil {
		log.Fatal("Error parsing environment variables.")
	}
}
