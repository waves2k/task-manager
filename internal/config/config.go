package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnectionString string
	Port             string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
		return nil, err
	}

	cfg := &Config{
		ConnectionString: os.Getenv("CONNECTION_STRING"),
		Port:             os.Getenv("PORT"),
	}
	return cfg, nil
}
