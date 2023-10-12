package main

/*
There is an idea that we can add some state thing and simply make a note a combination of a bericht and a day, but I feell ike
while that is possible we don't need that, and it would be more work than it would be worth
*/

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/mattn/go-sqlite3"
)


var (
	db *sql.DB 
	sqlite3Conn sqlite3.SQLiteConn

	personList = []string {"rick", "youri", "robert", "milan"}
	dayList = []string {"mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"}
	gridStates = []string {"E", "X", "O", "?"} 

	upgrader = websocket.Upgrader{}
	socketPool []*websocket.Conn

	noteList []Note;
	grid = loadGrid()

	berichte []BerichtQuery;

	updateGridStatement *sql.Stmt

	selectAllBerichte *sql.Stmt
	selectAllGrid *sql.Stmt
	selectAllNotes *sql.Stmt
)

type WSMessage struct {
	Command string `json:"command"`
	CurrentState string `json:"currentState"`
	Week string `json:"week"`
	Person string `json:"person"`
	Day string `json:"day"`
	OptID string `json:"optid"`
}

// TODO: actual logging, because that is useful

func main() {

	fmt.Println("starting...")
	
	var deployFlag = flag.Bool("deploy", false, "Enable the -deploy flag to actually deploy the server.")
	flag.Parse()

	var err error
	db, err = sql.Open("sqlite3", "file:./resources/DATABASE?cache=shared")

	defer db.Close()

	if err != nil {
		// TODO:
		log.Fatal("failed to open db with error ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping db with error ", err)
	}

	// creating prepared statements
	// insert bericht, remove bericht, change state
	updateGridStatement, err = db.Prepare(`UPDATE days AS d SET d.state = ? WHERE week = ? AND person = ? AND day = ?`)

	selectAllBerichte, err = db.Prepare(`SELECT * FROM berichte`) // TODO: would prob be fine not being a prepared statement
	selectAllGrid, err = db.Prepare(`SELECT * FROM days`)
	selectAllNotes, err = db.Prepare(`SELECT * FROM notes`)

	defer updateGridStatement.Close()
	defer selectAllBerichte.Close()
	defer selectAllGrid.Close()
	defer selectAllNotes.Close()

	berichte = GetAllBerichte(db)

	noteList = loadNotes(db)

	http.HandleFunc("/", MainRouteHandler)

	if *deployFlag {
		http.ListenAndServe(":80", nil)
		fmt.Println("started http server...")
	} else {
		http.ListenAndServe(":8000", nil)
		fmt.Println("started local http server...")
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
				bericht := BerichtQuery {
					0,
					m.CurrentState,
				}
				id := WriteNewBerichtGetId(db, bericht)
				m.OptID = strconv.Itoa(id)
				
				broadcast, err := json.Marshal(m)
				if err != nil {
					// TODO:
				}
				Broadcast(broadcast)

			} else if cmd == "del-bericht" {
				broadcast, err := json.Marshal(m)
				if err != nil {
					// TODO:
				}
	
				str, _ := strings.CutPrefix(m.OptID, "b")
				id, err := strconv.Atoi(str)
				if err != nil {
					// TODO:
				}
				RemoveBerichtById(db, id)
				Broadcast(broadcast)


			} else if cmd == "addnote" {
				newNote := Note {
					m.CurrentState,
					m.Week,
					m.Day,
					m.Person,
					"0",
				}
				// add the note to the db
				id := saveNoteGetID(db, newNote)

				// send this back to the client
				m.OptID = strconv.Itoa(id)
				broadcast, err := json.Marshal(m)
				if err != nil {
					// TODO:
				}
				Broadcast(broadcast)


			} else if cmd == "deletenote" {
				id, err := strconv.Atoi(m.CurrentState)
				if err != nil {
					// TODO:
				}
				removeNoteById(db, id)
				broadcast, err := json.Marshal(m)
				if err != nil {
					// TODO:
				}
				Broadcast(broadcast)


			} else if cmd == "open" {
				updateUserOnOpen(socketConn)
			}

			if err != nil {
				log.Println("Failed to read message with error", err)
				break
			}
		}
		saveGrid(grid)
	}
}

func updateUserOnOpen(socketConn *websocket.Conn) {

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

	noteList := loadNotes(db)

	for _, note := range noteList {
		// TODO: load the notes directly from the db
		n := WSMessage {
			Command: "addnote",
			CurrentState: note.Content,
			Week: note.Week,
			Person: note.Person,
			Day: note.Day,
			OptID: note.Id,
		}
		message, err := json.Marshal(n)

		if err != nil {
			// TODO
		}
		socketConn.WriteMessage(websocket.TextMessage, message)
	}

	berichte := GetAllBerichte(db)
	for _, bericht := range berichte {
		m := WSMessage {
			Command: "post-bericht",
			CurrentState: bericht.Content,
		}

		message, err := json.Marshal(m)

		if err != nil {
			// TODO:
		}

		socketConn.WriteMessage(websocket.TextMessage, message)
	}
}
