package main

import "errors"

type ShoppingItem struct {
	Id int `json:"id,string"`
	Content string `json:"content"`
	Action string `json:"action"`
}

func RemoveShoppingItemById(id int) error {
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
