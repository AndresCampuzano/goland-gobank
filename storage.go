package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Storage interface {
	CreateAccount(acc *Account) error
	UpdateAccount(acc *Account) error
	DeleteAccount(id string) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id string) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// Retrieve environment variables
	user := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	// Construct connection string
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS account (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            number BIGINT NOT NULL,
            balance BIGINT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
        INSERT INTO account (first_name, last_name, number, balance, created_at) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id
    `

	var id string
	err := s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted account
	acc.ID = id

	return nil
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id string) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var accounts []*Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id string) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %s not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
