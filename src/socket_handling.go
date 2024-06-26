package src

import (
	"golang.org/x/net/websocket"
	"github.com/2weiEmu/sis50.nl.go/pkg/lerror"
)

var WebSocketShopConnections []*websocket.Conn

func RemoveWebsocketFromPool(conn *websocket.Conn, list []*websocket.Conn) []*websocket.Conn {
	i := -1

	for j, ws := range list {
		if ws == conn {
			i = j
			break
		}
	}

	if i == -1 {
		lerror.ErrLog("There was no matching websocket connection found, ignoring.", nil)
		return nil
	}

	conn.Close()
	return append(list[:i], list[i+1:]...)
}

