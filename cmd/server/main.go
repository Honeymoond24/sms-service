package main

import (
	"database/sql"
	"github.com/Honeymoond24/sms-service/internal/config"
	"github.com/Honeymoond24/sms-service/internal/infrastructure/database"
	"log"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDB(cfg.DatabaseDSN)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error while closing database connection: %v", err)
		}
	}(db)
}
