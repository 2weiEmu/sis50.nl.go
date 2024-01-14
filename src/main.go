package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

type MessageStruct struct {
	Day string `json:"day"`
	Person string `json:"person"`
	State string `json:"state"`
}


var websocket_connections []*websocket.Conn
var p_ws_conn *string
var cal = ReadCalendar(InitCalendarDefault())

func main() {

	p_deploy := flag.Bool("d", false, "A flag specifying the deploy mode of the server.")
	p_port := flag.Int("p", 8000, "The port the server should be deployed on.")
	p_ws_conn = flag.String("ws", "ws://localhost:8000", "The websocket base")
	flag.Parse()

	router := mux.NewRouter()

	router.Handle("/dayWS", websocket.Handler(DayWebsocketHandler))
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
	index := "src/static/templates/index.html"
	tmpl, err := template.ParseFiles(index)
	if err != nil {
		fmt.Println(err)
	}

	err = tmpl.Execute(writer, p_ws_conn)
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
	http.ServeFile(writer, request, "src/static/templates/" + page + ".html")
}

func DayWebsocketHandler(conn *websocket.Conn) {
	fmt.Println("Activating WebSocket handler...")

	websocket_connections = append(websocket_connections, conn)
	fmt.Println(websocket_connections)

	var message MessageStruct
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
}

func BroadcastToConnections(message MessageStruct) {
	fmt.Println("[BROADCAST STARTING]")
	for i := 0; i < len(websocket_connections); i++ {
		fmt.Println("[WS] Sending to: ", websocket_connections[i])
		err := websocket.JSON.Send(websocket_connections[i], message)
		if err != nil {
			fmt.Println(err)
		}
	}
}


