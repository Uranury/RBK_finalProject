package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr     string
	RedisAddr      string
	DbURL          string
	MigrationsPath string
	JWTKey         string
	MailgunDomain  string
	MailgunAPIKey  string
}

type DBConfig struct {
	URL string `env:"DB_URL" required:"true"`
}

func Load() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	listenAddr := getEnv("LISTEN_ADDR", ":8080")
	redisAddr := getEnv("REDIS_ADDR", ":6379")
	dbURL := getEnv("DB_URL", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable")
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	JWTKey := os.Getenv("JWT_SECRET")
	MailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	MailgunAPIKey := os.Getenv("MAILGUN_API_KEY")

	if migrationsPath == "" {
		return nil, errors.New("migrations path not set")
	}

	if JWTKey == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	if MailgunDomain == "" || MailgunAPIKey == "" {
		log.Println("[WARN] Mailgun config not fully set – email features will be disabled")
	}

	return &Config{
		ListenAddr:     listenAddr,
		RedisAddr:      redisAddr,
		DbURL:          dbURL,
		MigrationsPath: migrationsPath,
		JWTKey:         JWTKey,
		MailgunDomain:  MailgunDomain,
		MailgunAPIKey:  MailgunAPIKey,
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
