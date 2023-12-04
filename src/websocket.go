package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Basically defining an enum
const (
	Present = iota
	PossiblyPresent
	PresentNoCook
	MaybeCooking
	Cooking
	NotPresent
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
		err := websocket.ReadJSON(&calendar_message);
		fmt.Println("Message received.")

		if err != nil {
			// TODO:
			fmt.Println("Closing websocket:", err)
			break
		}
		calendar_message.State = (calendar_message.State + 1) % 6
		// update the database

		websocket.WriteJSON(calendar_message)
	}
}

func RemoveWebsocketConnection(websocket *websocket.Conn) error {
	// finding the websocket by index
	index := -1
	for i := 0; i < len(websocket_list); i++ {
		if websocket_list[i] == websocket {
			index = i
		}
	}

	if index == -1 {
		// TODO:
		return errors.New("Failed to find a websocket to disconnect.")
	} else {
		websocket_list = append(websocket_list[:index], websocket_list[index + 1:]...)
		return nil
	}
}
