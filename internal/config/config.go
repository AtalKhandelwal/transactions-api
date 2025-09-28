package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	Port   string
}

func FromEnv() Config {
	return Config{
		DBHost: getenv("DB_HOST", "localhost"),
		DBPort: getenv("DB_PORT", "5432"),
		DBUser: getenv("DB_USER", "postgres"),
		DBPass: getenv("DB_PASSWORD", "postgres"),
		DBName: getenv("DB_NAME", "transactions"),
		Port:   getenv("PORT", "8080"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
