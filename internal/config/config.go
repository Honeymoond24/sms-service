package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseDSN string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Loading .env file failed: %v", err)
	}

	databaseDSN := os.Getenv("DATABASE_DSN")

	if databaseDSN == "" {
		log.Fatal("Environment variable is missing: DATABASE_DSN")
	}

	return &Config{
		DatabaseDSN: databaseDSN,
	}
}
