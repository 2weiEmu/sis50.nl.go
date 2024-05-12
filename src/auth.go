package main

import "net/http"

type UserAuth struct {

}

func NewUserAuthenticator() UserAuth {
	return UserAuth{}
}

// the idea is that the tokens will in the future maybe include more
// identifying information	
func WritePrivate(w http.ResponseWriter, cookie http.Cookie) {

}

func ReadPrivate(r *http.Request, name string) {

}
