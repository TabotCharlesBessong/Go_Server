package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var Envs = initconfig()

func initconfig() Config {
	godotenv.Load()
	return Config{
		DBHost: getENV("DB_HOST","localhost"),
		DBPort: getENV("DB_PORT","5432"),
		DBUser: getENV("DB_USER","postgres"),
		DBPassword: getENV("DB_PASSWORD","Charles#123"),
		DBName: getENV("DB_NAME","go_authentication"),
	}
}

func getENV(key,fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}