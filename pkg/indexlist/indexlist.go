package indexlist

import (
	"encoding/csv"
	"os"
	"strconv"

	c "sis50/pkg/constants"
	"sis50/pkg/lerror"
	"sis50/pkg/logger"
	n "sis50/pkg/node"
	s "sis50/pkg/shopping"
	"sis50/src"
	"github.com/BenLubar/htmlcleaner"
	"golang.org/x/net/websocket"
)

var ShopItemList, _ = ReadFromFile()
var IdCount = GetIdCount()


type IndexList struct {
	IndexList []n.IndexNode
}

func (list *IndexList) add(value s.ShoppingItem) {
	list.IndexList = append(list.IndexList, n.IndexNode {
		Index: list.Length(),
		Value: value,
	})
}

func GetIdCount() int {
	indexList, err := ReadFromFile()
	if err != nil {
		lerror.ErrLog("Failed reading from file when getting id count", err)
	}

	id := 0;

	for _, item := range indexList.IndexList {
		if item.Value.Id > id {
			id = item.Value.Id + 1
		}
	}

	return id
}


func (list *IndexList) Length() int {
	return len(list.IndexList)
}

func (list *IndexList) RemoveByItemId(id int) error {
	idx := list.IndexOfId(id)
	if idx == -1 {
		return lerror.ErrLog("Couldn't remove item", nil)
	}

	ridx := list.IndexList[idx].Index

	list.IndexList = append(list.IndexList[:idx], list.IndexList[idx+1:]...)

	for i, item := range list.IndexList {
		if item.Index > ridx {
			list.IndexList[i].Index -= 1
		}
	}
	return nil
}

func (list *IndexList) IndexOfId(id int) int {
	for i, item := range list.IndexList {
		if item.Value.Id == id {
			return i
		}
	}
	return -1
}

func (list *IndexList) EditMessageById(id int, newContent string) error {
	idx := list.IndexOfId(id)
	if idx == -1 {
		return lerror.ErrLog("Could not edit message using Id, id not found", nil)
	}

	list.IndexList[idx].Value.Content = newContent
	return nil
}

func (list *IndexList) MoveToNewIndexById(id int, newIndex int) error {
	idx := list.IndexOfId(id)
	oldIndex := list.IndexList[idx].Index

	if idx == -1 {
		return lerror.ErrLog("Id not found when moving", nil)
	}

	// if the old index was at the top of the list, we have to shift things up
	if oldIndex < newIndex {
		for i, item := range list.IndexList {
			if item.Index <= newIndex && item.Index >= oldIndex {
				list.IndexList[i].Index -= 1
			}
		}
	} else {
		for i, item := range list.IndexList {
			if item.Index >= newIndex && item.Index <= oldIndex {
				list.IndexList[i].Index += 1
			}
		}
	}

	list.IndexList[idx].Index = newIndex
	return nil
}

func (list *IndexList) WriteToFile() error {
	err := os.Truncate(c.ShoppingFile, 0)
	if err != nil {
		return lerror.ErrLog("Failed to truncate shopping file", err)
	}

	file, err := os.OpenFile(
		c.ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return lerror.ErrLog("Failed to open shopping file for writing", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range list.IndexList {
		err := writer.Write(item.Serialize())
		if err != nil {
			return lerror.ErrLog("Something went wrong writing to the file", err)
		}
	}

	writer.Flush()
	err = writer.Error()
	if err != nil {
		return lerror.ErrLog("The writer experienced an error when writing", err)
	}
	return nil
}

func NewIndexList() IndexList {
	return IndexList {
		make([]n.IndexNode, 0),
	}
}

func ReadFromFile() (IndexList, error) {
	file, err := os.OpenFile(
		c.ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return IndexList{}, 
			lerror.ErrLog(
				"Something went wrong when opening the file for reading", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return IndexList{}, lerror.ErrLog("Reader failed reading all records", err)
	}

	list := NewIndexList()

	for _, record := range records {
		deserialized, err := Deserialize(record)
		if err != nil {
			return IndexList{}, lerror.ErrLog("Failed to deserialize record", err)
		}

		list.IndexList = append(list.IndexList, deserialized)
	}

	return list, nil
}

func (list *IndexList) Ordered() []s.ShoppingItem {
	newList := make([]s.ShoppingItem, 0)

	for i := 0; i < list.Length(); i++ {
		for _, item := range list.IndexList {
			if item.Index == i {
				newList = append(newList, item.Value)
			}
		}
	}
	return newList
}

func Deserialize(serial []string) (n.IndexNode, error) {
	index, err := strconv.Atoi(serial[3])
	if err != nil {
		return n.IndexNode{}, lerror.ErrLog("Failed to convert index from file", err)
	}

	itemId, err := strconv.Atoi(serial[0])
	if err != nil {
		return n.IndexNode{}, lerror.ErrLog("Failed to convert item id from file", err)
	}

	return n.IndexNode {
		Index: index,
		Value: s.ShoppingItem{
			Id: itemId,
			Content: serial[1],
			Action: serial[2],
		},
	}, nil
}

func ShoppingListWebsocketHandler(conn *websocket.Conn) {
	logger.InfoLog.Println("Using Shopping Websocket Handler")
	src.WebSocketShopConnections = append(src.WebSocketShopConnections, conn)

	var message s.ShoppingItem

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

		keyword, prs := c.ShoppingActionMap[message.Action]

		if !prs {
			lerror.ErrLog("Action wasn't valid", nil)
			break
		}

		if keyword != c.OPEN {
			switch keyword {
			case c.ADD:
				message.Id = IdCount
				IdCount += 1
				ShopItemList.add(message)

			case c.REMOVE:
				err := ShopItemList.RemoveByItemId(message.Id)
				if err != nil {
					lerror.ErrLog("Failed to remove shopping item by id", err)
					break
				}
				
			case c.EDIT:
				err = ShopItemList.EditMessageById(message.Id, message.Content)
				if err != nil {
					lerror.ErrLog("Could not edit message by id", err)
					break
				}

			case c.REARRANGE:
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

			for _, wsConn := range src.WebSocketShopConnections {
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
	src.WebSocketShopConnections = n.RemoveWebsocketFromPool(
		conn, src.WebSocketShopConnections)
	conn.Close()
}
