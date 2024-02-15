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

var WebSocketDayConnections []*websocket.Conn
var WebSocketShopConnections []*websocket.Conn
var paramWebSocketConn *string
var StateCalendar = ReadCalendar(InitCalendarDefault())
var ShoppingList, err = ReadShoppingList()
var IdCount int
var AllMessagesList = readMessages(MessageList{});

const MessageFile = "./resources/messages"
const ShoppingFile = "./resources/shopping"
const CalendarFile = "./resources/calendar"

const DayCount = 7
const UserCount = 4

var StateList = []string{"present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"}
var ConstPersonList = []string{"rick", "youri", "robert", "milan"}
var DayList = []string{"ma", "di", "wo", "do", "vr", "za", "zo"}

func main() {
	paramDeploy := flag.Bool("d", false, "A flag specifying the deploy mode of the server.")
	paramPort := flag.Int("p", 8000, "The port the server should be deployed on.")
	paramWebSocketConn = flag.String("base", "localhost:8000", "The websocket base")
	flag.Parse()

	router := mux.NewRouter()

	router.Handle("/dayWS", websocket.Handler(DayWebsocketHandler))
	router.Handle("/shopWS", websocket.Handler(ShoppingListWebsocketHandler))
	router.HandleFunc("/api/messages/{pageNumber}", GETMessages).Methods("GET")
	router.HandleFunc("/api/messages", POSTMessage).Methods("POST")
	router.HandleFunc("/", IndexPage)
	router.HandleFunc("/{page}", GetPage)
	router.HandleFunc("/css/{style}", GetStyle)
	router.HandleFunc("/js/{script}", GetScript)
	router.HandleFunc("/images/{image}", GetImage)
	router.HandleFunc("/fonts/{font}", GetFont)
	router.HandleFunc("/admin", GetAdmin)
	http.Handle("/", router)

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
	pagesLength := len(AllMessagesList.Pages)
	if pagesLength == 0 {
		titleMsg = "No messages."
	} else {
		titleMsg = AllMessagesList.Pages[pagesLength - 1].Message[
			len(AllMessagesList.Pages[pagesLength - 1].Message) - 1]
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
	}
	
	err = tmpl.Execute(writer, paramWebSocketConn)
	if err != nil {
		fmt.Println(err)
	}
}

func getValidPages() []string {
	return []string{
		"messages", 
		"help",
	}
}
