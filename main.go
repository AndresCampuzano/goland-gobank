package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func seedAccount(s Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Andres", "CG", "password-1")
}

func main() {

	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		seedAccounts(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
