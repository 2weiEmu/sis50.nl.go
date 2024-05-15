package main

import (
	"database/sql"
	"net/http"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// i will open a db here, there is prob a better way to do this
	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
		
}

