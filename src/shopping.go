package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/benlubar/htmlcleaner"
	"golang.org/x/net/websocket"
)

type ShoppingItem struct {
	Id int `json:"id,string"`
	Content string `json:"content"`
	Action string `json:"action"`
}

func RemoveShoppingItemById(id int) error {
	fmt.Println("[INFO] Removing", id)
	i := -1

	for j, b := range ShoppingList {
		if b.Id == id {
			i = j
			break
		}
	}

	if i == -1 {
		// TODO:
		return errors.New("Failed to find the correct Id in the list")
	}

	ShoppingList = append(ShoppingList[:i], ShoppingList[i+1:]...)
	return nil
}

func EditMessageById(id int, content string) error {
	i := -1

	for j, b := range ShoppingList {
		if b.Id == id {
			i = j
			break
		}
	}

	if i == -1 {
		// TODO:
		return errors.New("Failed to find the correct Id in the list")
	}

	ShoppingList[i].Content = content
	return nil
}

func ReadShoppingList() ([]ShoppingItem, error) {
	file, err := os.OpenFile(
		ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	shoppingList := make([]ShoppingItem, 0)

	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		item := ShoppingItem{
			Id: id,
			Content: record[1],
			Action: record[2],
		}
		shoppingList = append(shoppingList, item)
	}

	fmt.Println("[INFO] Read:", shoppingList)
	return shoppingList, nil
}

func WriteShoppingList(shoppingList []ShoppingItem) error {
	err := os.Truncate(ShoppingFile, 0)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range shoppingList {
		serial := item.Serialize()
		fmt.Println("[INFO] Writing:", serial)
		writer.Write(serial)
	}

	writer.Flush()
	err = writer.Error()
	return err // will return nil if nil
}

func (item *ShoppingItem) Serialize() []string {
	list := make([]string, 3)

	list[0] = strconv.Itoa(item.Id)
	list[1] = item.Content
	list[2] = item.Action
	return list
}

func ShoppingListWebsocketHandler(shop_conn *websocket.Conn) {
	fmt.Println("[INFO] Activating Shopping Handler")

	WebSocketShopConnections = append(WebSocketShopConnections, shop_conn)
	fmt.Println(WebSocketShopConnections)

	var message ShoppingItem
	for {
		err := websocket.JSON.Receive(shop_conn, &message)
		if err != nil {
			// TODO:
			fmt.Println(err)
			break
		}

		fmt.Println("[INFO] Message received: ", message)
		message.Content = htmlcleaner.Clean(nil, message.Content)

		if message.Action != "open-shopping" {
			if message.Action == "remove" {
				err = RemoveShoppingItemById(message.Id)
				WriteShoppingList(ShoppingList)
				ShoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}

			} else if message.Action == "add" {
				message.Id = IdCount
				IdCount++
				ShoppingList = append(ShoppingList, message)
				WriteShoppingList(ShoppingList)
				ShoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}

			} else if message.Action == "edit" {
				err = EditMessageById(message.Id, message.Content)
				WriteShoppingList(ShoppingList)
				ShoppingList, err = ReadShoppingList()
				if err != nil {
					// TODO:
				}
			}

			for _, ws_conn := range WebSocketShopConnections {
				err = websocket.JSON.Send(ws_conn, message)
			}
			
			if err != nil {
				// TODO:
				fmt.Println(err)
			}
		} else {
			fmt.Println("Opening")
			for _, si := range ShoppingList {
				// NOTE: thought we had to make the actions be "add" manually -
				// but everything that gets added to the list already has "add"
				err := websocket.JSON.Send(shop_conn, si)
				if err != nil {
					// TODO:
					fmt.Println(err)
				}
			}
		}
	}
	WebSocketShopConnections = RemoveWebsocketFromPool(
		shop_conn, WebSocketShopConnections)
}
