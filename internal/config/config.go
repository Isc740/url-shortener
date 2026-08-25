package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func NewConfig() Config {
	port, err := getEnv("PORT")
	if err != nil {
		port = "8080"
	}

	dbUrl, err := getEnv("DATABASE_URL")
	if err != nil {
		dbUrl = "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"
	}

	return Config{
		Port:        port,
		DatabaseURL: dbUrl,
	}
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New("env variable not found")
	}
	return value, nil
}
