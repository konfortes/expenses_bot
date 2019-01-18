package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	fmt.Println("ronen")
	if err != nil {
		return nil, err
	}

	return db, nil
}
