package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/matthewhartstonge/argon2"
)

type LoginData struct {
	Username string `json:"username"`
	Given_password string `json:"password"`
}

func LoginUserPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginData := LoginData{}
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// now we have to check the password
	// if the password is right we give him a token and enter it in the db

	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Query("SELECT id, password_hash FROM users AS u WHERE u.username = ?", loginData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var encoded string
	var userId int
	if !result.Next() {
		http.Error(w, "dont seem to be in the db buddy", http.StatusInternalServerError)
		return
	}

	err = result.Scan(&userId, &encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result.Close()
	
	ok, err := argon2.VerifyEncoded([]byte(loginData.Given_password), []byte(encoded))
	
	if ok {
		sessionToken := MakeRandomString(64)
		sessionCookie := MakeSessionCookie(userId, string(sessionToken))
		_, err = db.Exec("INSERT INTO sessions (user_id, session_token) VALUES (? , ?)", userId, string(sessionToken))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = WritePrivate(w, sessionCookie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Incorrect Password", http.StatusUnauthorized)
}

func LogoutUserPost(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	DeleteCookie(w, "sis50session")
}

