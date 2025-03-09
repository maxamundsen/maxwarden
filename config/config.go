// The config package allows an easy way to provide variable options to the
// program read at runtime during startup. These are typically settings like the
// http port, database connection string, or password.
package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
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
