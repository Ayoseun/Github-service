package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// DatabaseConfig holds the configuration for the database connection.
type DatabaseConfig struct {
	DSN string
}

// Config holds the application configuration.
type Config struct {
	PER_PAGE          string `json:"PER_PAGE"`
	BASE_URL          string `json:"BASE_URL"`
	GITHUB_TOKEN      string `json:"GITHUB_TOKEN"`
	DEFAULT_OWNER     string `json:"DEFAULT_OWNER"`
	DEFAULT_REPO      string `json:"DEFAULT_REPO"`
	BEGIN_FETCH_DATE  string `json:"BEGIN_FETCH_DATE"`
	PORT              string `json:"PORT"`
	POSTGRES_USER     string `json:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `json:"POSTGRES_PASSWORD"`
	POSTGRES_HOST     string `json:"POSTGRES_HOST"`
	POSTGRES_DB       string `json:"POSTGRES_DB"`
	POLL_INTERVAL     int64  `json:"POLL_INTERVAL"`
}

// LoadConfig loads configuration from environment variables or a .env file.
// It first checks for a .env file and logs a message if it's not found.
// Then, it loads configuration from the environment variables or default values.
func LoadConfig() (config Config, err error) {
	// Set the config file path to .env
	viper.SetConfigFile(".env")
	// Allow reading from environment variables
	viper.AutomaticEnv()

	// Attempt to read the config file
	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		// Log a warning if the .env file is not found, but continue with env variables
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No .env file found, continuing with environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v\n", err)
			return
		}
	}

	// Unmarshal environment variables into the Config struct
	if err = viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshaling config: %v\n", err)
		return
	}

	log.Println("Configuration loaded successfully")
	return
}
