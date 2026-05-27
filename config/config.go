package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort        string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBDatabase        string
	DBSSLMode         string
	DBDriver          string
	PaymentServiceURL string
	PaymentProvider   string
	GCSBucket         string
	GCSPublicBaseURL  string
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string
	FrontAppURL       string
	AdminEmail        string
	AdminPassword     string
}

func Load() Config {
	return Config{
		ServerPort:        env("PORT", env("SERVER_PORT", "8085")),
		DBHost:            env("POSTGRES_HOST", "localhost"),
		DBPort:            env("POSTGRES_PORT", "5432"),
		DBUser:            env("POSTGRES_USER", "postgres"),
		DBPassword:        env("POSTGRES_PASSWORD", "postgres"),
		DBDatabase:        env("POSTGRES_DB", "adega"),
		DBSSLMode:         env("DB_SSL_MODE", "disable"),
		DBDriver:          env("DB_DRIVER", "postgres"),
		PaymentServiceURL: env("PAYMENT_SERVICE_URL", "http://localhost:8080/api/v1"),
		PaymentProvider:   env("PAYMENT_PROVIDER", "efi"),
		GCSBucket:         env("GCS_BUCKET", ""),
		GCSPublicBaseURL:  env("GCS_PUBLIC_BASE_URL", ""),
		SendGridAPIKey:    env("SENDGRID_API_KEY", ""),
		SendGridFromEmail: env("SENDGRID_FROM_EMAIL", ""),
		SendGridFromName:  env("SENDGRID_FROM_NAME", "Adega Flow"),
		FrontAppURL:       env("FRONT_APP_URL", "http://localhost:3000"),
		AdminEmail:        env("ADMIN_EMAIL", "admin@adega.com"),
		AdminPassword:     env("ADMIN_PASSWORD", "admin123"),
	}
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBDatabase,
		c.DBSSLMode,
	)
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
