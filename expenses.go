package main

import "fmt"

func persistExpense(amount float32, description string, userID int) {
	fmt.Printf("got an expense: amount: %f, description: %s\n", amount, description)

}
