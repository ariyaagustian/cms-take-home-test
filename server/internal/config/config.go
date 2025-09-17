package config

import (
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string

	JWTSecret string
	AppPort   string
	Env       string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		DBHost:    getenv("DB_HOST", "localhost"),
		DBPort:    getenv("DB_PORT", "5432"),
		DBUser:    getenv("DB_USER", "cms"),
		DBPass:    getenv("DB_PASS", "secret"),
		DBName:    getenv("DB_NAME", "cmsdb"),
		DBSSL:     getenv("DB_SSLMODE", "disable"),
		JWTSecret: getenv("JWT_SECRET", "change-me"),
		AppPort:   getenv("APP_PORT", "8080"),
		Env:       getenv("ENV", "development"),
	}
}
