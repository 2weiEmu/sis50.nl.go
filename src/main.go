package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var grid = loadGrid()

var personList = []string {"rick", "youri", "robert", "milan"}
var dayList = []string {"mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"}

var gridStates = []string {"E", "X", "O", "?"} 

var upgrader = websocket.Upgrader{}
var socketPool []*websocket.Conn

var noteList []Note = loadNotes()


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

func saveGrid(grid [][]string) {
	_ = os.Truncate("./src/resources/grid", 0)

	file, err := os.OpenFile("./src/resources/grid", os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println("failed to open save file")
	}

	defer file.Close()

	for _, arr := range grid {
		_, err := file.WriteString(strings.Join(arr, ",") + "\n")

		if err != nil {
			log.Fatal("failed writing line with error", err)
		}
	}

	fmt.Println("Finished saving grid to file.")

}

func loadGrid() [][]string {
	file, err := os.Open("./src/resources/grid")

	if err != nil {
		log.Fatal("failed to open grid file with error", err)
	}

	defer file.Close()

	var grid [][]string;

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		fmt.Println(line)

		grid = append(grid, line)
	}

	return grid;
}

func GetCmd(message string) string {
	return strings.Split(message, "$")[0]
}

func ParseToggle(message string) (week, person, day string) {
	arr := strings.Split(message, "$")
	return arr[1], arr[2], arr[3]
}

func ParseBericht(message string) string {
	return strings.Split(message, "$")[1]
}

func ModifyGrid(grid [][]string, week, person, day string) string {

	row := 0

	if week == "next" {
		row += 7
	}

	row += findIndex(dayList, day)

	col := findIndex(personList, person)

	grid[row][col] = GetNextMark(grid[row][col])

	return grid[row][col]
}

func GetNextMark(sy string) string {
	return gridStates[(findIndex(gridStates, sy) + 1) % len(gridStates)]
}

func findIndex(arr []string, s string) int {	

	for n, f := range arr {
		if f == s {
			return n;
		}
	}

	return -1
}

func Broadcast(message string) {

	for _, socket := range socketPool {
		socket.WriteMessage(websocket.TextMessage, []byte(message))
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
			_, message, err := socketConn.ReadMessage()

			cmd := GetCmd(string(message))

			if cmd == "toggle" {
				week, person, day := ParseToggle(string(message))
				fmt.Println(week, person, day)

				newState := ModifyGrid(grid, week, person, day)

				returnMessage := "toggle$" + newState + "$" + week + "$" + person + "$" + day

				// socketConn.WriteMessage(websocket.TextMessage, []byte(returnMessage))

				Broadcast(returnMessage)

			} else if cmd == "addnote" {
				addedNote := ParseNote(string(message))
				noteList = append(noteList, addedNote)

				Broadcast(string(message))
			} else if cmd == "deletenote" {
				deletedNote := ParseNote(string(message))

				removeIndex := -1

				for x, note := range noteList {
					if note == deletedNote {
						removeIndex = x 
						break
					}
				}

				noteList = RemoveIndex(noteList, removeIndex)

				Broadcast(string(message))

			} else if cmd == "post" {
				// TODO: make bericht do something
				content := ParseBericht(string(message))
				fmt.Println(content)

			} else if cmd == "open" {
				
				for i := 0; i < len(grid); i++ {
					for j := 0; j < len(grid[i]); j++ {

						week := "next"
						
						if i < 7 {
							week = "current"
						}
						
						person := personList[j]
						day := dayList[i % 7]

						message := "toggle$" + grid[i][j] + "$" + week + "$" + person + "$" + day

						socketConn.WriteMessage(websocket.TextMessage, []byte(message))
					}
				}

				for _, note := range noteList {

					message := "addnote$" + note.Content + "$" + note.Week + "$" + note.Person + "$" + note.Day

					socketConn.WriteMessage(websocket.TextMessage, []byte(message))
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
