package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

func LoadDatabaseConfig() (DatabaseConfig, error) {

	if err := godotenv.Load(); err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         parseEnvInt("DB_PORT", 3306),
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
	}, nil
}

func parseEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
