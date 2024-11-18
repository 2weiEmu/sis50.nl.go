package node

import (
	"strconv"

	"sis50/pkg/lerror"
	s "sis50/pkg/shopping"
	"golang.org/x/net/websocket"
)

type IndexNode struct {
	Index int;
	Value s.ShoppingItem;
}
func (node *IndexNode) Serialize() []string {
	return append(node.Value.Serialize(), strconv.Itoa(node.Index))
}

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

