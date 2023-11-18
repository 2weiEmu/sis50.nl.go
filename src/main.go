package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"crypto/rand"

	"github.com/gorilla/mux"
)

func main() {

	secretKey := make([]byte, 32)
	length, err := rand.Read(secretKey)

	if err != nil {
		// TODO:
	}
	
	if length != 32 {
		// TODO:
	}

	mux := mux.NewRouter().StrictSlash(true)

	mux.HandleFunc("/", GetMainPage)
	mux.HandleFunc("/css/{style}", GetCSSStyle)
	mux.HandleFunc("/js/{script}", GetJavaScript)
	mux.HandleFunc("/recipe/{recipe}", GetRecipeArticle)

	fmt.Println("Listening...")
	err = http.ListenAndServe(":8000", mux)

	if err != nil {
		fmt.Println("An error occured with the server:", err)
	}
}

// TODO: this _could_ be done better, for both the CSS and JS but explicit is good in this case
func GetCSSStyle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	style := vars["style"]
	
	http.ServeFile(writer, request, "src/static/css/" + style)
}

func GetJavaScript(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	script := vars["script"]

	http.ServeFile(writer, request, "src/static/js/" + script)
}

func GetMainPage(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "src/static/templates/index.html")
}

func ValidateLoginCookie(request *http.Request, name string, secret []byte, username string) bool {
	cookie, err := request.Cookie(name)

	if err != nil {
		// TODO:
		return false
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		// TODO:
		return false
	}

	// now we have the value of the cookie
	if len(value) < sha512.Size {
		// TODO:
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

	username := base64.URLEncoding.EncodeToString([]byte(user))
	cookieName := "login" + username

	cookie := http.Cookie {
		Name: cookieName,
		Value: username,
		Path: "/",
		MaxAge: 120, // NOTE: AGE OF 120 for TESTING
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	mac := hmac.New(sha512.New, secretKey)
	mac.Write([]byte(cookieName)) // the name of the cookie
	mac.Write([]byte(username)) // the value of the cookie
	hmac := mac.Sum(nil)

	cookie.Value = string(hmac) + cookie.Value

	return cookie

}

func GetRecipeArticle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	recipe := vars["recipe"]

	http.ServeFile(writer, request, "src/static/recipes/" + recipe)
}
