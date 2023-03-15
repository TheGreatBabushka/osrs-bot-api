package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type database_interface interface {
	getAllAccounts() ([]Account, error)

	getAccountByUsername(username string) (Account, error)
	getAccountByEmail(email string) (Account, error)
	getAccountByID(id int) (Account, error)
	putAccount(account Account) error
	deleteAccount(id int) error
	updateAccount(account Account) error
}

type database struct {
	db *sql.DB
}

func (d *database) getAllAccounts() ([]Account, error) {
	db := d.db

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT * FROM accounts")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		var email string
		var password string
		var status string

		if err := rows.Scan(&id, &username, &email, &password, &status); err != nil {
			log.Fatal(err)
		}

		log.Println(id, username, email, password, status)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return []Account{}, nil
}

func (d *database) getAccountByUsername(username string) (Account, error) {
	return Account{}, nil
}

func (d *database) getAccountByEmail(email string) (Account, error) {
	return Account{}, nil
}

func (d *database) getAccountByID(id int) (Account, error) {
	return Account{}, nil
}
