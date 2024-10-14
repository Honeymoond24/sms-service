package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	_ "modernc.org/sqlite"
	"time"
)

func NewDB(databaseDSN string) *sql.DB {
	db, err := sql.Open("sqlite", databaseDSN)

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
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
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		log.Fatalf("Failed to run migration: %v", err)
	}
}
