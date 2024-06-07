package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/2weiEmu/sis50.nl.go/src"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)


func main() {
	paramDeploy := flag.Bool(
		"d", false, "A flag specifying the deploy mode of the server.")
	paramPort := flag.Int(
		"p", 8000, "The port the server should be deployed on.")
	paramWebSocketConn := flag.String(
		"base", "localhost:8000", "Where websockets should connect.")
	cert := flag.String("c", "", "State the certificate location")
	secret := flag.String("k", "", "State the private key location")
	flag.Parse()

	logFile, err := os.OpenFile(src.MainLog, os.O_APPEND | os.O_RDWR, 664)
	if err != nil {
		fmt.Println("[ERROR] Failed to open main log file.")
		panic("Could not open log file.")
	}
	defer logFile.Close()

	src.InfoLog = log.New(logFile, "[INFO] ", src.LoggerFlags)
	src.RequestLog = log.New(logFile, "[REQUEST] ", src.LoggerFlags)
	src.ErrorLog = log.New(logFile, "[ERROR] ", src.LoggerFlags)

	cssDir := http.Dir("src/static/css")
	imgDir := http.Dir("src/static/images")
	jsDir := http.Dir("src/static/js")
	fontsDir := http.Dir("src/static/fonts")

	secure := "none"
	if *paramDeploy {
		secure = "ssl"
	} 

	HTMLctx, err := src.NewHTMLContext(
		src.LoggerFlags, logFile, secure, *paramWebSocketConn)
	if err != nil {
		fmt.Println("[ERROR] Could not create html context", err)
		panic("Nope, no context")
	}

	calHdl := src.NewCalendarHandler(src.LoggerFlags, logFile)

	router := mux.NewRouter()
	router.Handle("/dayWS",
		src.NewUserAuthenticator(websocket.Handler(calHdl.HandleCalendarWebsocket)))
	router.Handle("/shopWS",
		src.NewUserAuthenticator(websocket.Handler(src.ShoppingListWebsocketHandler)))

	router.Handle("/profile",
		src.NewUserAuthenticator(src.HandleFuncAsHandle(src.ReceiveUserProfileImage))).
		Methods("POST")
	router.HandleFunc("/api/messages/{pageNumber}", src.GETMessages).Methods("GET")
	router.HandleFunc("/api/messages", src.POSTMessage).Methods("POST")
	router.HandleFunc("/login", src.LoginUserPost).Methods("POST")
	router.Handle("/logout", 
		src.NewUserAuthenticator(src.HandleFuncAsHandle(src.LogoutUserPost))).
		Methods("POST")
	// this has to be wrapped with an auth because otherwise if you do it 
	// right you can log out arbitrary users	

	router.Handle("/", 
		src.NewUserAuthenticator(src.HandleFuncAsHandle(HTMLctx.HandleIndex)))
	router.HandleFunc("/login", HTMLctx.HandleLogin)
	router.Handle("/{page}",
		src.NewUserAuthenticator(src.HandleFuncAsHandle(HTMLctx.HandlePage)))
	// TODO: make above use file sever?

	http.Handle("/", router)

	// ok so apparently this doesn't work with router.Handle???
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(cssDir)))
	http.Handle("/images/",
		http.StripPrefix("/images/", http.FileServer(imgDir)))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(jsDir)))
	http.Handle("/fonts/", 
		http.StripPrefix("/fonts/", http.FileServer(fontsDir)))
	
	// go weeklyResetTimer()
	go calHdl.ShiftCalendarDaily()

	listenPort := ":" + strconv.Itoa(*paramPort)
	if *paramDeploy {
		fmt.Println("Began listening (SSL) on port:", listenPort);
		err = http.ListenAndServeTLS(listenPort, *cert, *secret, nil)

	} else {
		fmt.Println("Began listening on port:", listenPort);
		err = http.ListenAndServe(listenPort, nil)
	}

	if err != nil {
		src.ErrLog("Listen and serve failed with:", err)
	}
}

