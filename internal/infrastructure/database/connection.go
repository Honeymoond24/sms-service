package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	_ "modernc.org/sqlite"
)

func NewDB(databaseDSN string) *sql.DB {
	db, err := sql.Open("sqlite", databaseDSN)

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return db
}

func Migration(db *sql.DB) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/infrastructure/database/migrations",
		"sqlite", driver)
	if err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}
	err = m.Up()
	if err != nil {
		log.Printf("Failed to run migration: %v", err)
	}
}
