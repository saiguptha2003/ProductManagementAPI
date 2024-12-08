package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "./product_management.db")
	if err != nil {
		log.Fatalf("Failed to connect to the SQLite database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("SQLite database is unreachable: %v", err)
	}

	log.Println("SQLite database connection established")
}
