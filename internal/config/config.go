package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	SessionTTL  time.Duration
}

// LoadConfig loads config variables from environment variables
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	duration := os.Getenv("SESSION_TTL")
	if duration == "" {
		duration = "30s"
	}
	sessionTTL, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal("failed to parse session ttl: ", err)
	}

	return Config{
		Port:        port,
		DatabaseURL: databaseUrl,
		SessionTTL:  sessionTTL,
	}
}
