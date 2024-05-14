package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var conn sqlite3.SQLiteConn

// TODO: this is super temp and should be read from a config file
// that you set without any public knowledge
var secretKey []byte = []byte("")

type UserAuthWrapper struct {
	Db *sql.DB;
	handler http.Handler;
}

// method required to count as a handler itself
func (u *UserAuthWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := u.VerifySessionCookie(ReadPrivate(r, "sis50session"))
	if err != nil {
		http.Error(w, "Failed to verify the session", http.StatusUnauthorized)
		return
	}

	u.handler.ServeHTTP(w, r)
}

func NewUserAuthenticator() UserAuthWrapper {
	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		panic("Failed to open the db")
	}
	if err = db.Ping(); err != nil {
		panic("Failed to ping the database")
	}
	return UserAuthWrapper{
		Db: db,
	}
}

func (u *UserAuthWrapper) Close() {
	u.Db.Close()
}

func (u *UserAuthWrapper) VerifySessionCookie(value string) error {
	return errors.New("Failed to verify session cookie")
}

// the idea is that the tokens will in the future maybe include more
// identifying information	
func WritePrivate(w http.ResponseWriter, cookie http.Cookie) error {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	aes, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aes.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encrypted := aes.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(encrypted)

	return WriteCookie(w, cookie)
}

func WriteCookie(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return errors.New("Cookie value is too long")
	}

	http.SetCookie(w, &cookie)
	return nil
}

func ReadPrivate(r *http.Request, name string) (string, error) {
	encryptedVal, err := ReadCookie(r, name)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aes, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := aes.NonceSize()
	if len(encryptedVal) < nonceSize {
		return "", errors.New("Invalid value for the cookie")
	}

	nonce := encryptedVal[:nonceSize]
	ciphertext := encryptedVal[nonceSize:]

	text, err := aes.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", errors.New("failed to open value")
	}
	
	expected, val, ok := strings.Cut(string(text), ":")
	if !ok {
		return "", errors.New("invalid cookie value")
	}

	if expected != name {
		return "", errors.New("unexpected name match")
	}
	return val, nil
}


func ReadCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", nil
	}

	return string(value), nil
}
