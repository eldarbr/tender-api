package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress   string
	PostgresConnUrl string
	LogLevel        string
}

func GetEnv(key, defaultValue string, required bool) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	if required {
		return "", fmt.Errorf("required environment key %s is not set", key)
	}
	return defaultValue, nil
}

func processConfig(config *Config) error {
	serverAddress, err := GetEnv("SERVER_ADDRESS", ":8080", false)
	if err != nil {
		return err
	}
	config.ServerAddress = serverAddress

	postgresConnUrl, err := GetEnv("POSTGRES_CONN", "", true)
	if err != nil {
		return err
	}
	config.PostgresConnUrl = postgresConnUrl

	logLevel, err := GetEnv("LOG_LEVEL", "info", false)
	if err != nil {
		return err
	}
	config.LogLevel = logLevel

	return nil
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := processConfig(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
