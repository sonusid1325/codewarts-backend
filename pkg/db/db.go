package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes connection to the Postgres database and runs migrations.
func InitDB(host string, port int, user, password, dbname, sslmode string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var err error
	// Retry connection for database startup synchronization
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database connection... (Attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database after retries: %w", err)
	}

	log.Println("Successfully connected to the database.")

	if err := runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return DB, nil
}

func runMigrations() error {
	// Enable UUID extension
	_, err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Printf("Warning: uuid-ossp extension could not be enabled: %v. Using gen_random_uuid() or manual fallbacks if supported.", err)
	}

	// Create users table
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		is_verified BOOLEAN DEFAULT FALSE,
		verification_code VARCHAR(6),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(usersTableQuery); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// For backward compatibility on existing databases, alter the table if columns are missing
	// Setting default to TRUE for pre-existing records so they are not locked out
	if _, err := DB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT TRUE;"); err != nil {
		log.Printf("Migration Alter (is_verified) info/warning: %v", err)
	}
	if _, err := DB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_code VARCHAR(6);"); err != nil {
		log.Printf("Migration Alter (verification_code) info/warning: %v", err)
	}

	// Create user_progress table
	progressTableQuery := `
	CREATE TABLE IF NOT EXISTS user_progress (
		user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		current_chapter INT NOT NULL DEFAULT 1,
		current_task INT NOT NULL DEFAULT 1,
		completed_chapters JSONB NOT NULL DEFAULT '[]'::jsonb,
		container_id VARCHAR(255) DEFAULT '',
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(progressTableQuery); err != nil {
		return fmt.Errorf("failed to create user_progress table: %w", err)
	}

	log.Println("Database migrations completed successfully.")
	return nil
}
