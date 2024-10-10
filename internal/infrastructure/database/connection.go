package database

import (
	"database/sql"
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
