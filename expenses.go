package main

import (
	_ "github.com/lib/pq"
)

type Expense struct {
	Amount      float32
	Category    string
	UserID      int
	Description string
}

func persistExpense(expense *Expense) error {
	_, err := app.DB.Exec("INSERT INTO expenses(amount, category, user_id, description) VALUES($1, $2, $3, $4)", expense.Amount, expense.Category, expense.UserID, expense.Description)
	return err
}
