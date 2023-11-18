package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
)

func ValidateLoginCookie(request *http.Request, name string, secret []byte, username string) bool {
	cookie, err := request.Cookie(name)

	if err != nil {
		// TODO:
		fmt.Println("Failed to get cookie", err)
		return false
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		// TODO:
		fmt.Println("Failed to decode cookie", err)
		return false
	}

	// now we have the value of the cookie
	if len(value) < sha512.Size {
		// TODO:
		fmt.Println("Cookie too short")
		return false
	}

	signature := value[:sha512.Size]
	_ = value[sha512.Size:]

	// recalculate the HMAC
	mac := hmac.New(sha512.New, secret)
	mac.Write([]byte(name))
	mac.Write([]byte(username))
	expected := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expected) {
		// TODO:
		fmt.Println("Signatures and inequal")
		return false
	}

	// this also means the value checks out

	return true
}

/*
	This function generates a cookie validating the login of the given user
*/
func GenerateLoginCookie(user string, secretKey []byte) http.Cookie {
	/*
		In order to provide at least some resistance to people just using funky cookie editors
		we are going to attach a hash (together with a secretKey and also the name is included 
		generated on each server start, that not even I know)
		together with the original value of the cookie
	*/

	cookie := http.Cookie {
		Name: user,
		Value: user,
		Path: "/",
		MaxAge: 120, // NOTE: AGE OF 120 for TESTING
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	mac := hmac.New(sha512.New, secretKey)
	mac.Write([]byte(user)) // the name of the cookie
	mac.Write([]byte(user)) // the value of the cookie
	hmac := mac.Sum(nil)

	cookie.Value = base64.URLEncoding.EncodeToString([]byte(string(hmac) + user))

	return cookie

}


