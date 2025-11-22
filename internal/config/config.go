package config

import "os"

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	Database DatabaseConfig
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "haider"),
			Password: getEnv("DB_PASSWORD", "ha0770"),
			DBName:   getEnv("DB_NAME", "googleforms"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
