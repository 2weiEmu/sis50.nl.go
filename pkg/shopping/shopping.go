package shopping

import (
	"strconv"
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


