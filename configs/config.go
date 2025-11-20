package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	StaticDir string
	DbUrl     string
	MaxConns  int
}

var Envs = NewConfig()

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		return &Config{}
	}

	return &Config{
		Port:      getEnv("API_PORT", ":8080"),
		StaticDir: getEnv("STATIC_DIR", "./ui/html/"),
		DbUrl:     getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/snippetboxdb?sslmode=disable"),
		MaxConns:  getEnvAsInt("DB_MAX_CONNS", 10),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}
