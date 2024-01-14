package main

import "golang.org/x/net/websocket"

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
