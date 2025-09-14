package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	MONGODB_URI string
	DBName string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:   getEnv("PORT", "8080"),
		MONGODB_URI: getEnv("DB_URI", ""),
		DBName: getEnv("DB_NAME", "finpay"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
