package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
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

	return Config{
		Port:        port,
		DatabaseURL: databaseUrl,
	}
}
