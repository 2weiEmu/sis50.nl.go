package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var personList = []string {"rick", "youri", "robert", "milan"}
var dayList = []string {"mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"}
var gridStates = []string {"E", "X", "O", "?"} 

var upgrader = websocket.Upgrader{}
var socketPool []*websocket.Conn

var noteList []Note = loadNotes()
var grid = loadGrid()

type WSMessage struct {
	Command string `json:"command"`
	CurrentState string `json:"currentState"`
	Week string `json:"week"`
	Person string `json:"person"`
	Day string `json:"day"`
}

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

func findIndex(arr []string, s string) int {	

	for n, f := range arr {
		if f == s {
			return n;
		}
	}
	return -1
}

func Broadcast(message []byte) {

	for _, socket := range socketPool {
		socket.WriteMessage(websocket.TextMessage, message)
	}
}

func RemoveIndex(list []Note, index int) []Note {
    return append(list[:index], list[index+1:]...)
}


func MainRouteHandler(writer http.ResponseWriter, request *http.Request) {

	var path = request.URL.Path

	fmt.Println("path:", path)

	if path == "/" {
		http.Redirect(writer, request, "/koken", http.StatusSeeOther)
	} else if path == "/koken" || path == "/winkel" {
		http.ServeFile(writer, request, "./src/static/templates" + path + ".html")
	} else if path == "/js/koken.js" || path == "/js/koken-helper.js"{
		writer.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(writer, request, "./src/static" + path)
	} else if path == "/koken-ws" {

		socketConn, err := upgrader.Upgrade(writer, request, nil)

		socketPool = append(socketPool, socketConn)

		if err != nil {
			log.Print("Failed to create websocket connection with error", err)
		}

		defer socketConn.Close()


		for {
			var m WSMessage
			socketConn.ReadJSON(&m)
			fmt.Println(m)
			cmd := m.Command	

			if cmd == "toggle" {
				newState := ModifyGrid(grid, m.Week, m.Person, m.Day)
				m.CurrentState = newState
				returnM, err := json.Marshal(m)

				if err != nil {
					// TODO:
				}
				Broadcast(returnM)

			} else if cmd == "post-bericht" {
				// TODO: make bericht do something
				//content := ParseBericht(string(message))
				//fmt.Println(content)


			} else if cmd == "open" {
				for i := 0; i < len(grid); i++ {
					for j := 0; j < len(grid[i]); j++ {

						week := "next"
						
						if i < 7 {
							week = "current"
						}
						person := personList[j]
						day := dayList[i % 7]

						n := WSMessage {
							Command: "toggle",
							CurrentState: grid[i][j],
							Week: week,
							Person: person,
							Day: day,
						}
						message, err := json.Marshal(n)
						if err != nil {
							// TODO
						}
						socketConn.WriteMessage(websocket.TextMessage, message)
					}
				}

				for _, note := range noteList {
					n := WSMessage {
						Command: "addnote",
						CurrentState: note.Content,
						Week: note.Week,
						Person: note.Person,
						Day: note.Day,
					}
					message, err := json.Marshal(n)

					if err != nil {
						// TODO
					}
					socketConn.WriteMessage(websocket.TextMessage, message)
				}
			}

			if err != nil {
				log.Println("Failed to read message with error", err)
				break
			}
		}
		saveNotes(noteList)
		saveGrid(grid)
	}
}
