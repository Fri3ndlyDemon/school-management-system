package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"

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

func RunMigrations(db *sql.DB, migrationsDir string) {
	// read all files in the migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// collect and sort the .sql file names to ensure correct order
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// execute each migration file sequentially
	for _, fileName := range sqlFiles {
		fullPath := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", fileName, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", fileName, err)
		}

		log.Printf("Successfully applied migration: %s", fileName)
	}

	log.Println("All database migrations applied successfully!")
}
