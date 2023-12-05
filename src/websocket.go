package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Basically defining an enum
const (
	Present = iota		// 0
	PossiblyPresent		// 1
	PresentNoCook		// 2
	MaybeCooking		// 3
	Cooking				// 4
	NotPresent			// 5
)

type CalendarMessage struct {
	Person string `json:"person"`
	Day string `json:"day"`
	State int `json:"calendar_state"`
}

func CalendarWebsocket(writer http.ResponseWriter, request *http.Request) {
	websocket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		// TODO:
	}
	defer websocket.Close()


	websocket_list = append(websocket_list, websocket)

	// Reading and writing messages

	var calendar_message CalendarMessage
	for {
		err := websocket.ReadJSON(&calendar_message)
		if err != nil {
			// TODO:
			fmt.Println("Closing websocket:", err)
			break
		}

		// update the database
		calendar_message.State = (calendar_message.State + 1) % 6

		// broadcast the message
		BroadcastMessageJSON(calendar_message)
	}

	RemoveWebsocketConnection(websocket)
}

func RemoveWebsocketConnection(websocket *websocket.Conn) error {
	index := -1
	for i := 0; i < len(websocket_list); i++ {
		if websocket_list[i] == websocket {
			index = i
		}
	}

	if index == -1 {
		return errors.New("Failed to find a websocket to disconnect.")
	} else {
		websocket_list = append(websocket_list[:index], websocket_list[index + 1:]...)
		return nil
	}
}

func BroadcastMessageJSON(message interface{}) {
	for i := 0; i < len(websocket_list); i++ {
		websocket_list[i].WriteJSON(message)
	}
}
