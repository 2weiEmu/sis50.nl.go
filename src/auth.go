package main

import (
	"database/sql"
	"net/http"

	"github.com/mattn/go-sqlite3"
)

var conn sqlite3.SQLiteConn

type UserAuth struct {
	Db *sql.DB;
}

func NewUserAuthenticator() UserAuth {
	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		panic("Failed to open the db")
	}
	if err = db.Ping(); err != nil {
		panic("Failed to ping the database")
	}
	return UserAuth{
		Db: db,
	}
}

func (u *UserAuth) Close() {
	u.Db.Close()
}

func (u *UserAuth) VerifySessionCookie(cookie http.Cookie) {
	
}

// the idea is that the tokens will in the future maybe include more
// identifying information	
func WritePrivate(w http.ResponseWriter, cookie http.Cookie) {

}

func ReadPrivate(r *http.Request, name string) {

}
