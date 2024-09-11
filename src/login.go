package src

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/matthewhartstonge/argon2"
	"sis50.nl.go/pkg/auth"
)

type LoginData struct {
	Username string `json:"username"`
	Given_password string `json:"password"`
}

func LoginUserPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "You are not in the database", http.StatusInternalServerError)
		return
	}

	loginData := LoginData{}
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		http.Error(w, "You are not in the database", http.StatusInternalServerError)
		return
	}

	// now we have to check the password
	// if the password is right we give him a token and enter it in the db

	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		http.Error(w, "You are not in the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Query("SELECT id, password_hash FROM users AS u WHERE u.username = ?", loginData.Username)
	if err != nil {
		http.Error(w, "You are not in the database", http.StatusInternalServerError)
		return
	}

	var encoded string
	var userId int
	if !result.Next() {
		http.Error(w, "You are not in the database", http.StatusInternalServerError)
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
		sessionCookie := auth.MakeSessionCookie(userId, string(sessionToken))
		_, err = db.Exec("INSERT INTO sessions (user_id, session_token) VALUES (? , ?)", userId, string(sessionToken))
		if err != nil {
			// im stupid we dont need the templates in the login post handler
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = auth.WritePrivate(w, sessionCookie)
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
	defer db.Close()

	userId, err := auth.GetUserIdFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM sessions AS s WHERE s.user_id = ?", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// i dont think we have to delete the cookie but its probably good practice
	auth.DeleteCookie(w, "sis50session")
	w.WriteHeader(http.StatusOK)
}

