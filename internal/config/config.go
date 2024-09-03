package config

import (
	"log"

	"github.com/spf13/viper"
)

// DatabaseConfig struct holds the configuration for the database connection
type DatabaseConfig struct {
	DSN string
}

// Config struct holds the application configuration
type Config struct {
	DATABASE_DEV_URL  string `json:"DATABASE_DEV_URL"`
	DATABASE_PROD_URL string `json:"DATABASE_PROD_URL"`
	BASE_URL          string `json:"BASE_URL"`
	GITHUB_TOKEN      string `json:"GITHUB_TOKEN"`
}

// LoadConfig loads configuration from environment variables or a file
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// GetDatabaseConfig returns the database configuration based on the environment
func GetDatabaseConfig(cfg Config, env string) DatabaseConfig {
	var dsn string

	// Determine the DSN based on the environment
	switch env {
	case "dev":
		dsn = cfg.DATABASE_DEV_URL
	case "prod":
		dsn = cfg.DATABASE_PROD_URL
	default:
		log.Fatalf("Unknown environment: %s", env)
	}

	return DatabaseConfig{
		DSN: dsn,
	}
}
