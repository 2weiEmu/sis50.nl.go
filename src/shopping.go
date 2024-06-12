package src

import (
	"strconv"

	"github.com/benlubar/htmlcleaner"
	"github.com/2weiEmu/sis50.nl.go/pkg/lerror"
	"github.com/2weiEmu/sis50.nl.go/pkg/logger"
	"golang.org/x/net/websocket"
)

var ShopItemList, _ = ReadFromFile()
var IdCount = GetIdCount()

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

func GetIdCount() int {
	indexList, err := ReadFromFile()
	if err != nil {
		lerror.ErrLog("Failed reading from file when getting id count", err)
	}

	id := 0;

	for _, item := range indexList.indexList {
		if item.value.Id > id {
			id = item.value.Id + 1
		}
	}

	return id
}

func ShoppingListWebsocketHandler(conn *websocket.Conn) {
	logger.InfoLog.Println("Using Shopping Websocket Handler")
	webSocketShopConnections = append(webSocketShopConnections, conn)

	var message ShoppingItem

	for {
		err := websocket.JSON.Receive(conn, &message)
		if err != nil {
			lerror.ErrLog("Failed to read websocket JSON", err)
			break
		}

		logger.InfoLog.Println("Message received:", message)
		ShopItemList, err = ReadFromFile()
		if err != nil {
			lerror.ErrLog("Failed to re-read shopping list from file", err)
			break
		}

		message.Content = htmlcleaner.Clean(nil, message.Content)

		keyword, prs := ShoppingActionMap[message.Action]

		if !prs {
			lerror.ErrLog("Action wasn't valid", nil)
			break
		}

		if keyword != OPEN {
			switch keyword {
			case ADD:
				message.Id = IdCount
				IdCount += 1
				ShopItemList.add(message)

			case REMOVE:
				err := ShopItemList.RemoveByItemId(message.Id)
				if err != nil {
					lerror.ErrLog("Failed to remove shopping item by id", err)
					break
				}
				
			case EDIT:
				err = ShopItemList.EditMessageById(message.Id, message.Content)
				if err != nil {
					lerror.ErrLog("Could not edit message by id", err)
					break
				}

			case REARRANGE:
				newIdx, err := strconv.Atoi(message.Content)
				if err != nil {
					lerror.ErrLog("Failed to convert message content when rearranging", err)
					break
				}

				err = ShopItemList.MoveToNewIndexById(message.Id, newIdx)
				if err != nil {
					lerror.ErrLog("Failed to move to new index", err)
					break
				}
			}

			err = ShopItemList.WriteToFile()
			if err != nil {
				lerror.ErrLog("Failed to write shopping list to file", err)
				break
			}

			for _, wsConn := range webSocketShopConnections {
				err = websocket.JSON.Send(wsConn, message)
			}
			if err != nil {
				lerror.ErrLog("Failed to broadcast to other connections", err)
				break
			}

		} else {
			logger.InfoLog.Println("Sending new opening websocket")

			for _, item := range ShopItemList.Ordered() {
				err := websocket.JSON.Send(conn, item)
				if err != nil {
					lerror.ErrLog("Failed to send opening shopping list statement", err)
					break
				}
			}
			logger.InfoLog.Println("Completed sending opening")
		}

	}
	webSocketShopConnections = RemoveWebsocketFromPool(
		conn, webSocketShopConnections)
	conn.Close()
}

