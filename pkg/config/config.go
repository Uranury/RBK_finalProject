package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ListenAddr     string
	RedisAddr      string
	DbURL          string
	MigrationsPath string
}

func Load() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	listenAddr := getEnv("LISTEN_ADDR", ":8080")
	redisAddr := getEnv("REDIS_ADDR", ":6379")
	dbURL := getEnv("DB_URL", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable")
	migrationsPath := os.Getenv("MIGRATIONS_PATH")

	if migrationsPath == "" {
		return nil, errors.New("migrations path not set")
	}

	return &Config{
		ListenAddr:     listenAddr,
		RedisAddr:      redisAddr,
		DbURL:          dbURL,
		MigrationsPath: migrationsPath,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadEnv() error {
	return godotenv.Load()
}
