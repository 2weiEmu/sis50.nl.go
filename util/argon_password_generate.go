package main

import (
	"database/sql"

	"github.com/matthewhartstonge/argon2"
	"github.com/mattn/go-sqlite3"
)

func main() {
	// the idea is to only generate an account, with a salt
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte("password"))
    if err != nil {
        panic("Failed to hash the password " + err.Error()) 
    }

	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		panic("Failed to open the database " + err.Error())
	}
	defer db.Close()

	var _ sqlite3.SQLiteConn

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ('test', ?)", encoded)
	if err != nil {
		panic("Failed to insert the values into the database " + err.Error())
	}
}
