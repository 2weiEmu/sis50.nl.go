package main

import (
	"flag"
	"fmt"
	"net/http"
)


func main() {
	

	var deployFlag = flag.Bool("deploy", false, "Enable the -deploy flag to actually deploy the server.")

	flag.Parse()

	http.HandleFunc("/", MainRouteHandler)

	if *deployFlag {
		http.ListenAndServe(":80", nil)
	} else {
		http.ListenAndServe(":8000", nil)
	}
}

func MainRouteHandler(writer http.ResponseWriter, request *http.Request) {

	var path = request.URL.Path

	fmt.Println("path:", path)

	if path == "/" {
		http.Redirect(writer, request, "/koken", http.StatusSeeOther)
	} else if path == "/koken" || path == "/winkel" {
		http.ServeFile(writer, request, "./src/static/templates" + path + ".html")
	} else if path == "/js/koken.js" {
		http.ServeFile(writer, request, "./src/static/js/koken.js")
	}

}
