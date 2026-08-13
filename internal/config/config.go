package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const minSecretLen = 32

type Config struct {
	DatabaseDSN   string
	JWTSecret     []byte
	JWTExpiration time.Duration
	JWTIssuer     string
	ServerAddr    string
}

type loader struct {
	errs []error
}

func (l *loader) required(key string) string {
	value := os.Getenv(key)
	if value == "" {
		l.errs = append(l.errs, fmt.Errorf("%s is required", key))
	}
	return value
}

func (l *loader) optional(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var l loader

	secret := l.required("JWT_SECRET")
	if secret != "" && len(secret) < minSecretLen {
		l.errs = append(l.errs, fmt.Errorf("JWT_SECRET must be at least %d bytes, but got %d", minSecretLen, len(secret)))
	}

	cfg := &Config{
		DatabaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			l.required("DB_HOST"),
			l.required("DB_USER"),
			l.required("DB_PASSWORD"),
			l.required("DB_NAME"),
			l.required("DB_PORT"),
			l.optional("SSL_MODE", "disable"),
		),
		JWTSecret:     []byte(secret),
		// Curto de propósito: o access token não é revogável. A renovação é
		// responsabilidade do refresh token, que vive no banco e pode ser revogado.
		JWTExpiration: 15 * time.Minute,
		JWTIssuer:     l.optional("JWT_ISSUER", "authsys"),
		ServerAddr:    ":" + l.optional("PORT", "8080"),
	}

	if len(l.errs) > 0 {
		return nil, errors.Join(l.errs...)
	}

	return cfg, nil
}
