package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseDSN string
	PushSMSURL  string
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

	pushSMSURL := os.Getenv("PUSH_SMS_URL")
	if pushSMSURL == "" {
		log.Fatal("Environment variable is missing: PUSH_SMS_URL")
	}

	return &Config{
		DatabaseDSN: databaseDSN,
		PushSMSURL:  pushSMSURL,
	}
}
