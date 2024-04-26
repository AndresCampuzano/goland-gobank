package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccount   int `json:"to_account"`
	FromAccount int `json:"from_account"`
	Amount      int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string) *Account {
	// Generate a random number between 10000000 and 99999999 (inclusive)
	// Ensures that the number does not start or end with zero
	number := int64(rand.Intn(90000000) + 10000000)

	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    number,
		CreatedAt: time.Now().UTC(),
	}
}
