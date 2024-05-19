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
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

var conn sqlite3.SQLiteConn

// TODO: this is super temp and should be read from a config file
// that you set without any public knowledge
var secretKey []byte = []byte("gb+V6PcZ8PC7oObI/kngTjBHrYsNKQ==")

type UserAuthWrapper struct {
	Db *sql.DB;
	TokenStmt *sql.Stmt
	handler http.Handler;
}

// method required to count as a handler itself
func (u UserAuthWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionval, err := ReadPrivate(r, "sis50session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = u.VerifySessionToken(sessionval)
	if err != nil {
		http.Error(w, "Failed to verify the session" + err.Error(), http.StatusUnauthorized)
		return
	}

	u.handler.ServeHTTP(w, r)
}

func NewUserAuthenticator(handler http.Handler) UserAuthWrapper {
	db, err := sql.Open("sqlite3", "./resources/centralDb")
	if err != nil {
		panic("Failed to open the db")
	}
	tokenstmt, err := db.Prepare("SELECT session_token FROM sessions AS s WHERE s.user_id =	?")
	if err != nil {
		panic("Failed to set up prepared statement")
	}
	return UserAuthWrapper{
		Db: db,
		TokenStmt: tokenstmt,
		handler: handler,
	}
}

func (u *UserAuthWrapper) Close() {
	u.TokenStmt.Close()
	u.Db.Close()
}

func MakeSessionCookie(userId int, token string) http.Cookie {
	return http.Cookie{
		Name: "sis50session",
		Value: strconv.Itoa(userId) + "$" + token,
		Path: "/",
		MaxAge: 2629800, // apparently this may not be great???
		HttpOnly: true, // cannot be modified with javascript
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (u *UserAuthWrapper) VerifySessionToken(value string) error {
	vals := strings.Split(value, "$")
	userId, err := strconv.Atoi(vals[0])
	if err != nil {
		return err
	}

	givenToken := vals[1]
	// TODO: limit max amount of logins at the same time with different tokens
	// TODO: use transactions
	if err = u.Db.Ping(); err != nil {
		return err
	}
	rows, err := u.TokenStmt.Query(userId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		err := rows.Scan(&token)
		fmt.Println(token, givenToken)
		if err != nil {
			return err
		}
		
		if token == givenToken {
			return nil
		}
	}
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
	fmt.Println("wrote cookie")
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

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name: name,
		Value: "",
		Path: "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true, // cannot be modified with javascript
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	})
}
