package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// connect database with databaseURL
func ConnectDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	// ping test to check
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully!")
	return db, nil
}
