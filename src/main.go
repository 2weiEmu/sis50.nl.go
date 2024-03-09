package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var webSocketDayConnections []*websocket.Conn
var webSocketShopConnections []*websocket.Conn
var paramWebSocketConn *string
var stateCalendar = ReadCalendar(InitCalendarDefault())
var shoppingList, err = ReadShoppingList()
var idCount int
var allMessagesList = readMessages(MessageList{});

func main() {
	paramDeploy := flag.Bool(
		"d", false, "A flag specifying the deploy mode of the server.")
	paramPort := flag.Int(
		"p", 8000, "The port the server should be deployed on.")
	paramWebSocketConn = flag.String(
		"base", "localhost:8000", "Where websockets should connect.")
	flag.Parse()

	// NOTE: consider changing this, to something like src/static and moving
	// all the templates somewhere else so they can't be accessed... maybe a
	// "public" folder?
	cssDir := http.Dir("src/static/css")
	imgDir := http.Dir("src/static/images")
	jsDir := http.Dir("src/static/js")
	fontsDir := http.Dir("src/static/fonts")

	router := mux.NewRouter()
	router.Handle("/dayWS", websocket.Handler(DayWebsocketHandler))
	router.Handle("/shopWS", websocket.Handler(ShoppingListWebsocketHandler))
	router.HandleFunc("/api/messages/{pageNumber}", GETMessages).Methods("GET")
	router.HandleFunc("/api/messages", POSTMessage).Methods("POST")
	router.HandleFunc("/", IndexPage)
	router.HandleFunc("/{page}", GetPage)
	router.HandleFunc("/admin", GetAdmin)
	http.Handle("/", router)

	// ok so apparently this doesn't work with router.Handle???
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(cssDir)))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(imgDir)))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(jsDir)))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(fontsDir)))
	
	// go weeklyResetTimer()
	go shiftCalendarDaily()

	if *paramDeploy {
		// TODO:
		// http.ListenAndServeTLS()
	} else {

		fmt.Println("Began listening on port: " + strconv.Itoa(*paramPort));
		http.ListenAndServe(":" + strconv.Itoa(*paramPort), nil)
	}
}

type IndexPageStruct struct {
	Message string
	Args string
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	index := "src/static/templates/index.html"
	tmpl, err := template.ParseFiles(index)
	if err != nil {
		fmt.Println(err)
	}

	var titleMsg string
	pagesLength := len(allMessagesList.Pages)
	if pagesLength == 0 {
		titleMsg = "No messages."
	} else {
		titleMsg = allMessagesList.Pages[pagesLength - 1].Message[
			len(allMessagesList.Pages[pagesLength - 1].Message) - 1]
	}

	MainPageStruct := IndexPageStruct{
		Message: titleMsg,
		Args: *paramWebSocketConn,
	}

	err = tmpl.Execute(writer, MainPageStruct)
	if err != nil {
		fmt.Println(err)
	}
}

func GetAdmin(writer http.ResponseWriter, request *http.Request) {
}

func GetPage(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	page := vars["page"]

	if !slices.Contains(getValidPages(), page) {
		http.ServeFile(
			writer, request, "src/static/templates/404.html",
		)
		return
	}

	pageLocation := "src/static/templates/" + page + ".html"
	tmpl, err := template.ParseFiles(pageLocation)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(
			writer, request, "src/static/templates/500.html",
		)
		return
	}
	
	err = tmpl.Execute(writer, paramWebSocketConn)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(
			writer, request, "src/static/templates/500.html",
		)
	}
}
