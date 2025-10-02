package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB(filepath string) {
	var err error
	DB, err = sqlx.Connect("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	schema := `
    CREATE TABLE IF NOT EXISTS products (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        quantity INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP
    );
    `
	DB.MustExec(schema)
}
