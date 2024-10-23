package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

func NewDB(databaseDSN string) *sql.DB {
	db, err := sql.Open("pgx", databaseDSN)

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

func Migration(databaseDSN string) {
	m, err := migrate.New(
		"file://./internal/infrastructure/database/migrations",
		databaseDSN)
	if err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		log.Fatalf("Failed to run migration: %v", err)
	}
}
