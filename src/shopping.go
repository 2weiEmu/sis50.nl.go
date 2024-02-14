package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type ShoppingItem struct {
	Id int `json:"id,string"`
	Content string `json:"content"`
	Action string `json:"action"`
}

func RemoveShoppingItemById(id int) error {
	fmt.Println("[INFO] Removing", id)
	i := -1

	for j, b := range shoppingList {
		if b.Id == id {
			i = j
			break
		}
	}

	if i == -1 {
		// TODO:
		return errors.New("Failed to find the correct Id in the list")
	}

	shoppingList = append(shoppingList[:i], shoppingList[i+1:]...)
	return nil
}

func EditMessageById(id int, content string) error {
	i := -1

	for j, b := range shoppingList {
		if b.Id == id {
			i = j
			break
		}
	}

	if i == -1 {
		// TODO:
		return errors.New("Failed to find the correct Id in the list")
	}

	shoppingList[i].Content = content
	return nil
}

func ReadShoppingList() ([]ShoppingItem, error) {
	file, err := os.OpenFile(
		SHOPPING_FILE, os.O_RDWR | os.O_APPEND, os.ModeAppend)
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
	err := os.Truncate(SHOPPING_FILE, 0)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		SHOPPING_FILE, os.O_RDWR | os.O_APPEND, os.ModeAppend)
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
