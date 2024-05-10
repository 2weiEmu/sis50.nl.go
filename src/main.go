package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var webSocketDayConnections []*websocket.Conn
var webSocketShopConnections []*websocket.Conn
var stateCalendar = ReadCalendar(InitCalendarDefault())
var shopItemList, _ = ReadFromFile()
var idCount = getIdCount()
var allMessagesList, _ = readMessages(MessageList{});
var infoLog, requestLog, errorLog *log.Logger

func main() {

	var secure string
	paramDeploy := flag.Bool(
		"d", false, "A flag specifying the deploy mode of the server.")
	paramPort := flag.Int(
		"p", 8000, "The port the server should be deployed on.")
	paramWebSocketConn := flag.String(
		"base", "localhost:8000", "Where websockets should connect.")
	cert := flag.String("c", "", "State the certificate location")
	secret := flag.String("k", "", "State the private key location")
	flag.Parse()

	logFile, err := os.OpenFile("./log/sis50.log", os.O_APPEND | os.O_RDWR, 664)
	if err != nil {
		fmt.Println("[LOGS] Failed to open main log file.")
		os.Exit(-1)
	}
	defer logFile.Close()

	infoLog = log.New(logFile, "[INFO] ", loggerFlags)
	requestLog = log.New(logFile, "[REQUEST] ", loggerFlags)
	errorLog = log.New(logFile, "[ERROR] ", loggerFlags)

	// NOTE: consider changing this, to something like src/static and moving
	// all the templates somewhere else so they can't be accessed... maybe a
	// "public" folder?
	cssDir := http.Dir("src/static/css")
	imgDir := http.Dir("src/static/images")
	jsDir := http.Dir("src/static/js")
	fontsDir := http.Dir("src/static/fonts")

	if *paramDeploy {
		secure = "ssl"
	} else {
		secure = "none"
	}

	HTMLctx, err := initHTMLContext(loggerFlags, logFile, secure, *paramWebSocketConn)

	router := mux.NewRouter()
	router.Handle("/dayWS", websocket.Handler(DayWebsocketHandler))
	router.Handle("/shopWS", websocket.Handler(ShoppingListWebsocketHandler))
	router.HandleFunc("/api/messages/{pageNumber}", GETMessages).Methods("GET")
	router.HandleFunc("/api/messages", POSTMessage).Methods("POST")
	router.HandleFunc("/", HTMLctx.HandleIndex)
	router.HandleFunc("/{page}", HTMLctx.HandlePage)
	http.Handle("/", router)

	// ok so apparently this doesn't work with router.Handle???
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(cssDir)))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(imgDir)))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(jsDir)))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(fontsDir)))
	
	// go weeklyResetTimer()
	go shiftCalendarDaily()

	if *paramDeploy {
		secure = "ssl"
		fmt.Println("Began listening (SSL) on port: " + strconv.Itoa(*paramPort));
		err = http.ListenAndServeTLS(":" + strconv.Itoa(*paramPort), *cert, *secret, nil)

	} else {
		secure = "none"
		fmt.Println("Began listening on port: " + strconv.Itoa(*paramPort));
		err = http.ListenAndServe(":" + strconv.Itoa(*paramPort), nil)
	}

	if err != nil {
		ErrLog("Listen and serve failed with:", err)
	}
}

type IndexPageStruct struct {
	Message string
	Args string
}

