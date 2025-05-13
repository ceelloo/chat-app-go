package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func New(addr, maxIdleTime string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	var version string
	db.QueryRow("SELECT sqlite_version()").Scan(&version)
	
	log.Printf("SQLite Version: %s", version)
	log.Printf("Connected to SQLite: %s", addr)
	return db, nil
}

func LoadMigration(dir string, db *sql.DB) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read migrations directory", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("Failed to read migration file", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal("Failed to apply migration", err)
		}

		log.Printf("Applied migration: %s", path)
	}

	return nil
}
