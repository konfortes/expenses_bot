package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func persistExpense(amount float32, description string, userID int) {
	fmt.Printf("got an expense: amount: %f, description: %s\n", amount, description)
	app.DB.Query("INSERT INTO expenses(amount, category, user_id, description) VALUES($1, $2, $3, $4)", amount, "category", userID, description)
}
