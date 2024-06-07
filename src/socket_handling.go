package src

import (
	"golang.org/x/net/websocket"
)

var webSocketShopConnections []*websocket.Conn

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

	conn.Close()
	return append(list[:i], list[i+1:]...)
}

func (handler *CalendarHandler) BroadcastToConnections(message CalMessage) {
	handler.InfoLog.Println("Broadcasting Websocket")
	for i := 0; i < len(handler.Connections); i++ {
		InfoLog.Println("Sending to: ", handler.Connections[i])
		err := websocket.JSON.Send(handler.Connections[i], message)
		if err != nil {
			ErrLog("Failed to send JSON via websocket during broadcast", err)
		}
	}
}
