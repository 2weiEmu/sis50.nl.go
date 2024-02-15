package main

import (
	"fmt"

	"golang.org/x/net/websocket"
)

func RemoveWebsocketFromPool(conn *websocket.Conn, list []*websocket.Conn) []*websocket.Conn {
	i := -1

	for j, ws := range list {
		if ws == conn {
			i = j
			break
		}
	}

	if i == -1 {
		// TODO:
		return nil
	}
	
	return append(list[:i], list[i+1:]...)
}

func BroadcastToConnections(message CalMessage) {
	fmt.Println("[BROADCAST STARTING]")
	for i := 0; i < len(WebSocketDayConnections); i++ {
		fmt.Println("[WS] Sending to: ", WebSocketDayConnections[i])
		err := websocket.JSON.Send(WebSocketDayConnections[i], message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
