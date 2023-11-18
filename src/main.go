package main

import (
	"fmt"
	"net/http"
	"crypto/rand"

	"github.com/gorilla/mux"
)

var secretKey []byte;

var testUserName = "temp"

func main() {

	secretKey = make([]byte, 32)
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
	mux.HandleFunc("/cookies/handout-test", handoutCookieTest)
	mux.HandleFunc("/cookies/validate-test", validateCookieTest)

	fmt.Println("Listening...")
	err = http.ListenAndServe(":8000", mux)

	if err != nil {
		fmt.Println("An error occured with the server:", err)
	}
}

func handoutCookieTest(writer http.ResponseWriter, request *http.Request) {
	cookie := GenerateLoginCookie("temp", secretKey)

	http.SetCookie(writer, &cookie)

	fmt.Println("Cookie Set")
	writer.Write([]byte("Set the cookie " + (cookie.Name) + " for the user temp. The cookie value is as follows:" + string(cookie.Value) + "\n"))
}

func validateCookieTest(writer http.ResponseWriter, request *http.Request) {
	if val, user := ValidateLoginCookie(request, "sis50.nl.login.validation", secretKey, "temp"); val {
		fmt.Print("===================\nCookie was validated for user:", user ,"\n===================\n")
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

func GetRecipeArticle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	recipe := vars["recipe"]

	http.ServeFile(writer, request, "src/static/recipes/" + recipe)
}
