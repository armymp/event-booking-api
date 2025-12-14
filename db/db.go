package db

import (
	"database/sql"
	"log"

	"github.com/armymp/event-booking-api/config"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	dbName := config.AppConfig.Database.Name
	if dbName == "" {
		dbName = "events.db" // fallback
	}

	log.Printf("Using database file: %s (Environment: %s)\n", dbName, config.AppConfig.Server.Env)

	var err error
	DB, err = sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("Failed to enable foreign key enforcement: %v", err)
	}

	// Optional: connection pool settings
	DB.SetMaxOpenConns(0)
	DB.SetMaxIdleConns(2)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Could not create events table: %v", err)
	}
}