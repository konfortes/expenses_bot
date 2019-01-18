package main

import (
	_ "github.com/lib/pq"
)

func persistExpense(amount float32, description string, userID int) error {
	_, err := app.DB.Exec("INSERT INTO expenses(amount, category, user_id, description) VALUES($1, $2, $3, $4)", amount, "category", userID, description)
	return err
}
