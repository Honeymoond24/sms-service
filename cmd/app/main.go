package main

import (
	"database/sql"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/Honeymoond24/sms-service/internal/config"
	"github.com/Honeymoond24/sms-service/internal/infrastructure/database"
	"github.com/Honeymoond24/sms-service/internal/interfaces/rest"
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

	database.Migration(db)

	servicesRepository := database.NewServicesRepository(db)
	smsService := application.NewSmsService(servicesRepository)

	rest.NewServer(smsService).Run()
}
