package main

import (
	"fmt"
	"strconv"

	"github.com/benlubar/htmlcleaner"
	"golang.org/x/net/websocket"
)

type ShoppingItem struct {
	Id int `json:"id,string"`
	Content string `json:"content"`
	Action string `json:"action"`
}

func (item *ShoppingItem) Serialize() []string {
	list := make([]string, 3)

	list[0] = strconv.Itoa(item.Id)
	list[1] = item.Content
	list[2] = item.Action
	return list
}

func ShoppingListWebsocketHandler(conn *websocket.Conn) {
	infoLog.Println("Using Shopping Websocket Handler")
	webSocketShopConnections = append(webSocketShopConnections, conn)

	var message ShoppingItem

	for {
		err := websocket.JSON.Receive(conn, &message)
		if err != nil {
			ErrLog("Failed to read websocket JSON", err)
			fmt.Println("Failed with message received:", err)
			break
		}

		infoLog.Println("Message received:", message)
		shopItemList, err = ReadFromFile()
		if err != nil {
			ErrLog("Failed to re-read shopping list from file", err)
			break
		}

		message.Content = htmlcleaner.Clean(nil, message.Content)

		keyword, prs := shoppingActionMap[message.Action]

		if !prs {
			ErrLog("Action wasn't valid", nil)
			break
		}

		if keyword != OPEN {
			switch keyword {
			case ADD:
				message.Id = idCount
				idCount++
				shopItemList.add(message)

			case REMOVE:
				err := shopItemList.RemoveByItemId(message.Id)
				if err != nil {
					ErrLog("Failed to remove shopping item by id", err)
					break
				}
				
			case EDIT:
				err = shopItemList.EditMessageById(message.Id, message.Content)
				if err != nil {
					ErrLog("Could not edit message by id", err)
					break
				}

			case REARRANGE:
			}

			err = shopItemList.WriteToFile()
			if err != nil {
				ErrLog("Failed to write shopping list to file", err)
				break
			}

			for _, wsConn := range webSocketShopConnections {
				err = websocket.JSON.Send(wsConn, message)
			}
			if err != nil {
				ErrLog("Failed to broadcast to other connections", err)
				break
			}

		} else {
			infoLog.Println("Sending new opening websocket")

			for _, node := range shopItemList.indexList {
				err := websocket.JSON.Send(conn, node.value)
				if err != nil {
					ErrLog("Failed to send opening shopping list statement", err)
					break
				}
			}
			infoLog.Println("Completed sending opening")
		}

	}
	webSocketShopConnections = RemoveWebsocketFromPool(
		conn, webSocketShopConnections)
	conn.Close()
}

