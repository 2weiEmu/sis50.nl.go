package main

import (
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
		ErrLog("There was no matching websocket connection found, ignoring.", nil)
		return nil
	}
	
	return append(list[:i], list[i+1:]...)
}

func BroadcastToConnections(message CalMessage) {
	infoLog.Println("Broadcasting Websocket")
	for i := 0; i < len(webSocketDayConnections); i++ {
		infoLog.Println("Sending to: ", webSocketDayConnections[i])
		err := websocket.JSON.Send(webSocketDayConnections[i], message)
		if err != nil {
			ErrLog("Failed to send JSON via websocket during broadcast", err)
		}
	}
}
