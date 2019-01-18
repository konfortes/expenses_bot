package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
