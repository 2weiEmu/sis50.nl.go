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

var upgrader = websocket.Upgrader{}

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

	file, err := os.Open("./src/resources/grid")

	if err != nil {
		fmt.Println("failed to open save file")
	}

	for _, arr := range grid {
		_, err := file.WriteString(strings.Join(arr, "") + "\n")

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


func MainRouteHandler(writer http.ResponseWriter, request *http.Request) {

	var path = request.URL.Path

	fmt.Println("path:", path)

	if path == "/" {
		http.Redirect(writer, request, "/koken", http.StatusSeeOther)
	} else if path == "/koken" || path == "/winkel" {
		http.ServeFile(writer, request, "./src/static/templates" + path + ".html")
	} else if path == "/js/koken.js" {
		http.ServeFile(writer, request, "./src/static/js/koken.js")
	} else if path == "/koken-ws" {

		socketConn, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Print("Failed to create websocket connection with error", err)
		}

		defer socketConn.Close()


		for {
			mt, message, err := socketConn.ReadMessage()

			cmd := GetCmd(string(message))

			if cmd == "toggle" {
				week, person, day := ParseToggle(string(message))
				fmt.Println(week, person, day)

			} else if cmd == "post" {
				content := ParseBericht(string(message))
				fmt.Println(content)

			}

			if err != nil {
				log.Println("Failed to read message with error", err)
				break
			}

			

			err = socketConn.WriteMessage(mt, []byte("hi"))
		}

		saveGrid(grid)

	}

}
