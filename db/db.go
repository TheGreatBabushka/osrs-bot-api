package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Represents a row in the accounts table
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

// Represents a row in the levels table
type Levels struct {
	Attack       int `json:"attack"`
	Strength     int `json:"strength"`
	Defence      int `json:"defence"`
	Ranged       int `json:"ranged"`
	Magic        int `json:"magic"`
	Prayer       int `json:"prayer"`
	Runecrafting int `json:"runecrafting"`
	Hitpoints    int `json:"hitpoints"`
	Agility      int `json:"agility"`
	Herblore     int `json:"herblore"`
	Thieving     int `json:"thieving"`
	Crafting     int `json:"crafting"`
	Fletching    int `json:"fletching"`
	Slayer       int `json:"slayer"`
	Hunter       int `json:"hunter"`
	Mining       int `json:"mining"`
	Smithing     int `json:"smithing"`
	Fishing      int `json:"fishing"`
	Cooking      int `json:"cooking"`
	Firemaking   int `json:"firemaking"`
	Woodcutting  int `json:"woodcutting"`
	Farming      int `json:"farming"`
}

type Activity struct {
	AccountID int    `json:"account_id"`
	Command   string `json:"command"`
	StartedAt string `json:"started_at"`
}

type Database struct {
	Driver *sql.DB
}

func (d *Database) GetAccounts() ([]Account, error) {
	db := d.Driver
	accounts := []Account{}

	stmtOut, err := db.Prepare("SELECT * FROM accounts")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		log.Fatal(err)
		return accounts, err
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
		// log.Println(id, username, email, password, status)

		accounts = append(accounts, Account{
			ID:       id,
			Username: username,
			Email:    email,
			Password: password,
			Status:   status,
		})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return accounts, err
	}

	return accounts, nil
}

func (d *Database) GetAccount(id int) (Account, error) {
	db := d.Driver
	var account Account

	stmtOut, err := db.Prepare("SELECT * FROM accounts WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	row := stmtOut.QueryRow(id)
	err = row.Scan(&account.ID, &account.Username, &account.Email, &account.Password, &account.Status)
	if err != nil {
		log.Fatal(err)
		return account, err
	}

	return account, nil
}

func (d *Database) InsertAccount(email string, username string) {
	db := d.Driver

	stmtOut, err := db.Prepare("INSERT INTO accounts (email, username, status, password) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE email = ?, username = ?, status = ?, password = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtOut.Exec(email, username, "active", "password", email, username, "active", "password")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Account Inserted")
}

func (d *Database) GetAccountByEmail(email string) (Account, error) {
	db := d.Driver
	var account Account

	stmtOut, err := db.Prepare("SELECT * FROM accounts WHERE email = ?")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer stmtOut.Close()

	row := stmtOut.QueryRow(email)
	err = row.Scan(&account.ID, &account.Username, &account.Email, &account.Password, &account.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			// insert the account

		}
		log.Println(err)
		return account, err
	}

	return account, nil
}

func (d *Database) GetLevelsForAccount(id int) (Levels, error) {
	db := d.Driver

	rows, err := db.Query("SELECT * FROM levels WHERE account_id = ?", id)
	if err != nil {
		log.Fatal(err)
		return Levels{}, err
	}
	defer rows.Close()

	lvls := Levels{}
	for rows.Next() {
		if err := rows.Scan(&lvls.Attack, &lvls.Strength, &lvls.Defence, &lvls.Ranged, &lvls.Magic, &lvls.Prayer, &lvls.Runecrafting, &lvls.Hitpoints, &lvls.Agility, &lvls.Herblore, &lvls.Thieving, &lvls.Crafting, &lvls.Fletching, &lvls.Slayer, &lvls.Hunter, &lvls.Mining, &lvls.Smithing, &lvls.Fishing, &lvls.Cooking, &lvls.Firemaking, &lvls.Woodcutting, &lvls.Farming); err != nil {
			log.Fatal(err)
			return Levels{}, err
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return Levels{}, err
	}

	// fmt.Println(lvls)
	return lvls, nil
}

func (d *Database) UpdateLevelsForAccount(acc Account, lvls Levels) error {
	// get the names of the columns in the levels table
	columns, err := d.LevelsColumns()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// fmt.Println(columns)

	// build the query - the ON DUPLICATE KEY UPDATE part is what makes this work (UPSERT)
	query := "INSERT INTO levels (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s"
	var columnsString string
	var valuesString string
	var updateString string

	for i, column := range columns {
		if i == 0 {
			columnsString += column
			valuesString += "?"
			updateString += column + " = ?"
		} else {
			columnsString += ", " + column
			valuesString += ", ?"
			updateString += ", " + column + " = ?"
		}
	}
	query = fmt.Sprintf(query, columnsString, valuesString, updateString)
	// fmt.Println(query)

	values := []interface{}{}
	values = append(values, acc.ID, lvls.Attack, lvls.Strength, lvls.Defence, lvls.Ranged, lvls.Magic, lvls.Prayer, lvls.Runecrafting, lvls.Hitpoints, lvls.Agility, lvls.Herblore, lvls.Thieving, lvls.Crafting, lvls.Fletching, lvls.Slayer, lvls.Hunter, lvls.Mining, lvls.Smithing, lvls.Fishing, lvls.Cooking, lvls.Firemaking, lvls.Woodcutting, lvls.Farming)
	values = append(values, acc.ID, lvls.Attack, lvls.Strength, lvls.Defence, lvls.Ranged, lvls.Magic, lvls.Prayer, lvls.Runecrafting, lvls.Hitpoints, lvls.Agility, lvls.Herblore, lvls.Thieving, lvls.Crafting, lvls.Fletching, lvls.Slayer, lvls.Hunter, lvls.Mining, lvls.Smithing, lvls.Fishing, lvls.Cooking, lvls.Firemaking, lvls.Woodcutting, lvls.Farming)

	err = d.prepareExecute(query, values...)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (d *Database) LevelsColumns() ([]string, error) {
	rows, err := d.prepareQuery("SELECT * FROM levels LIMIT 1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return columns, nil
}

func (d *Database) Query(query string) (*sql.Rows, error) {
	return d.prepareQuery(query)
}

func (d *Database) prepareExecute(query string, args ...interface{}) error {
	db := d.Driver

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (d *Database) prepareQuery(query string, args ...interface{}) (*sql.Rows, error) {
	db := d.Driver

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return rows, nil
}