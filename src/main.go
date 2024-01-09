package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)



func main() {

	p_deploy := flag.Bool("d", false, "A flag specifying the deploy mode of the server.")
	p_port := flag.Int("p", 8000, "The port the server should be deployed on.")
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/", IndexPage)
	router.HandleFunc("/{page}", GetPage)
	router.HandleFunc("/css/{style}", GetStyle)
	router.HandleFunc("/js/{script}", GetScript)
	router.HandleFunc("/images/{image}", GetImage)
	router.HandleFunc("/fonts/{font}", GetFont)
	router.HandleFunc("/admin", GetAdmin)
	http.Handle("/", router)

	if *p_deploy {
		// TODO:
		// http.ListenAndServeTLS()
	} else {

		fmt.Println("Began listening on port: " + strconv.Itoa(*p_port));
		http.ListenAndServe(":" + strconv.Itoa(*p_port), nil)
	}
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer ,request, "src/static/templates/index.html")
}

func GetStyle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	style := vars["style"]
	http.ServeFile(writer, request, "src/static/css/" + style)
}

func GetScript(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	script := vars["script"]
	http.ServeFile(writer, request, "src/static/js/" + script)
}

func GetImage(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	image := vars["image"]
	http.ServeFile(writer, request, "src/static/images/" + image)
}

func GetFont(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	font := vars["font"]
	http.ServeFile(writer, request, "src/static/fonts/" + font)
}

func GetAdmin(writer http.ResponseWriter, request *http.Request) {
}

func GetPage(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	page := vars["page"]
	http.ServeFile(writer, request, "src/static/templates/" + page + ".html")
}
