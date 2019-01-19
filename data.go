package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connectedChannel := make(chan bool, 1)
	var db *sql.DB
	var err error

	go func() {
		connStr := os.Getenv("DB_URL")
		for {
			db, err = sql.Open("postgres", connStr)
			if err == nil {
				err = db.Ping()
				if err == nil {
					connectedChannel <- true
				}
			}
		}
	}()

	select {
	case <-connectedChannel:
		return db, nil
	case <-time.After(2 * time.Second):
		return nil, err
	}
}
