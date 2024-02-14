package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/benlubar/htmlcleaner"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var websocket_day_connections []*websocket.Conn
var websocket_shop_connections []*websocket.Conn
var p_ws_conn *string
var cal = ReadCalendar(InitCalendarDefault())
var shoppingList, err = ReadShoppingList()
var id_count int

var messageList = readMessages(MessageList{});

const SHOPPING_FILE = "./resources/shopping"

func main() {
	p_deploy := flag.Bool("d", false, "A flag specifying the deploy mode of the server.")
	p_port := flag.Int("p", 8000, "The port the server should be deployed on.")
	p_ws_conn = flag.String("base", "localhost:8000", "The websocket base")
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

	if *p_deploy {
		// TODO:
		// http.ListenAndServeTLS()
	} else {

		fmt.Println("Began listening on port: " + strconv.Itoa(*p_port));
		http.ListenAndServe(":" + strconv.Itoa(*p_port), nil)
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
	pagesLength := len(messageList.Pages)
	if pagesLength == 0 {
		titleMsg = "No messages."
	} else {
		titleMsg = messageList.Pages[pagesLength - 1].Message[
			len(messageList.Pages[pagesLength - 1].Message) - 1]
	}

	MainPageStruct := IndexPageStruct{
		Message: titleMsg,
		Args: *p_ws_conn,
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

	pageLocation := "src/static/templates/" + page + ".html"
	tmpl, err := template.ParseFiles(pageLocation)
	if err != nil {
		fmt.Println(err)
	}
	
	err = tmpl.Execute(writer, p_ws_conn)
	if err != nil {
		fmt.Println(err)
	}
}

func DayWebsocketHandler(conn *websocket.Conn) {
	fmt.Println("Activating WebSocket handler...")

	websocket_day_connections = append(websocket_day_connections, conn)
	fmt.Println(websocket_day_connections)

	var message CalMessage
	for {
		err := websocket.JSON.Receive(conn, &message)
		if err != nil {
			// TODO:
			fmt.Println(err)
			break
		}
		fmt.Println("Message received: ", message)

		if message.State != "open-calendar" {
			message.State = UpdateCalendar(cal, message)
			BroadcastToConnections(message)
		} else {
			m := ""
			for _, s := range cal.Day {
				for _, k := range s {
					m += strconv.Itoa(k)
				}
				m += "/"
			}
			fmt.Println("[INFO] Open:", m)

			message.Day = m
			err := websocket.JSON.Send(conn, &message)
			if err != nil {
				fmt.Println(err)
				// TODO:
			}
		}
	}
	WriteCalendar(cal)
	websocket_day_connections = RemoveWebsocketFromPool(conn, websocket_day_connections)
}


func ShoppingListWebsocketHandler(shop_conn *websocket.Conn) {
	fmt.Println("[INFO] Activating Shopping Handler")

	websocket_shop_connections = append(websocket_shop_connections, shop_conn)
	fmt.Println(websocket_shop_connections)

	var message ShoppingItem
	for {
		err := websocket.JSON.Receive(shop_conn, &message)
		if err != nil {
			// TODO:
			fmt.Println(err)
			break
		}

		fmt.Println("[INFO] Message received: ", message)
		message.Content = htmlcleaner.Clean(nil, message.Content)

		if message.Action != "open-shopping" {
			if message.Action == "remove" {
				err = RemoveShoppingItemById(message.Id)
				WriteShoppingList(shoppingList)
				shoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}

			} else if message.Action == "add" {
				message.Id = id_count
				id_count++
				shoppingList = append(shoppingList, message)
				WriteShoppingList(shoppingList)
				shoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}

			} else if message.Action == "edit" {
				err = EditMessageById(message.Id, message.Content)
				WriteShoppingList(shoppingList)
				shoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}
			}

			for _, ws_conn := range websocket_shop_connections {
				err = websocket.JSON.Send(ws_conn, message)
			}
			
			if err != nil {
				// TODO:
				fmt.Println(err)
			}
		} else {
			fmt.Println("Opening")
			for _, si := range shoppingList {
				// NOTE: thought we had to make the actions be "add" manually -
				// but everything that gets added to the list already has "add"
				err := websocket.JSON.Send(shop_conn, si)
				if err != nil {
					// TODO:
					fmt.Println(err)
				}
			}
		}
	}
	websocket_shop_connections = RemoveWebsocketFromPool(
		shop_conn, websocket_shop_connections)
}

func BroadcastToConnections(message CalMessage) {
	fmt.Println("[BROADCAST STARTING]")
	for i := 0; i < len(websocket_day_connections); i++ {
		fmt.Println("[WS] Sending to: ", websocket_day_connections[i])
		err := websocket.JSON.Send(websocket_day_connections[i], message)
		if err != nil {
			fmt.Println(err)
		}
	}
}


