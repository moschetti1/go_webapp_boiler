package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	PublicHost     string
	DBURL          string
	CookieSecret   string
	GoogleClientId string
	GoogleSecret   string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()
		cfg = &Config{
			Port:           getEnv("PORT", ":8000"),
			PublicHost:     getEnv("HOST", "http://localhost"),
			DBURL:          getEnv("GOOSE_DBSTRING", "file:local.db"),
			CookieSecret:   getEnv("COOKIE_SECRET", "s3cr37"),
			GoogleClientId: getEnv("GOOGLE_CLIENT_ID", "asdf"),
			GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET", "asdf"),
		}
	})
	return cfg
}

func MustGet() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded: call config.Load() first")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvOrError(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s is not set", key))

}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return b
	}

	return fallback
}
