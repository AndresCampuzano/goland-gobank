package main

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type LoginResponse struct {
	Number int    `json:"number"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	Number   int    `json:"number"`
	Password string `json:"password"`
}

// TransferRequest FIXME: finish
type TransferRequest struct {
	ToAccount   int `json:"to_account"`
	FromAccount int `json:"from_account"`
	Amount      int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int       `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int       `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw))
	if err != nil {
		return false
	}

	return true
}

func NewAccount(firstName string, lastName string, password string) (*Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Generate a random number between 10000000 and 99999999 (inclusive)
	// Ensures that the number does not start or end with zero
	number := rand.Intn(90000000) + 10000000

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            number,
		EncryptedPassword: string(encryptedPassword),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
