package main

import (
	"encoding/json"
	"io"
	"net/http"
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
}

