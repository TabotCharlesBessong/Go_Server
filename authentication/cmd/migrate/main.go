package main

import (
	"authentication/config"
	"authentication/db"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// Initialize database connection
	dbConn, err := db.NewPostgreSQL(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}
	defer dbConn.Close()

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(dbConn); err != nil {
		log.Fatal("error creating migrations table: ", err)
	}

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		log.Fatal("error getting migration files: ", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(dbConn)
	if err != nil {
		log.Fatal("error getting applied migrations: ", err)
	}

	// Apply new migrations
	for _, migration := range migrations {
		if !applied[migration] {
			if err := applyMigration(dbConn, migration); err != nil {
				log.Fatal("error applying migration: ", err)
			}
			log.Printf("Applied migration: %s\n", migration)
		}
	}

	log.Println("Migrations completed successfully")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}
	return applied, nil
}

func applyMigration(db *sql.DB, migration string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join("migrations", migration))
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO migrations (name) VALUES ($1)", migration); err != nil {
		return err
	}

	return tx.Commit()
}
