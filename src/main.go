package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
    fmt.Println("Starting the server...")

    deployFlag := flag.Bool("deploy", false, "Set the deployment mode with -deploy.")

    flag.Parse()

    http.HandleFunc("/", MainRouteHandler)

    var err error

    if !*deployFlag {
        err = http.ListenAndServe(":8000", nil)

    } else {
        err = http.ListenAndServe(":80", nil)

    }

    if err != nil {
        log.Fatal("Something crashed with err...", err)
    }
}

func MainRouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path	

	if requestPath == "/" {

		http.Redirect(writer, request, "/koken", http.StatusSeeOther)
	} else if requestPath == "/koken" || requestPath == "/bericht" || requestPath == "/comic" || requestPath == "/winkel" {

		http.ServeFile(writer, request, "src/static/templates" + requestPath + ".html")

		fmt.Println("Serving src/static/templates" + requestPath + ".html")
	} else if requestPath == "/css/index.css" {
		http.ServeFile(writer, request, "src/static/css/index.css")
		fmt.Println("Serving index.css")

	}
}
