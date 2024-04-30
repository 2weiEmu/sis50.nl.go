package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type IndexNode struct {
	index int;
	value ShoppingItem;
}

func (node *IndexNode) Serialize() []string {
	return append(node.value.Serialize(), strconv.Itoa(node.index))
}

type IndexList struct {
	indexList []IndexNode
}

func (list *IndexList) add(value ShoppingItem) {
	list.indexList = append(list.indexList, IndexNode {
		index: len(list.indexList),
		value: value,
	})
}

func (list *IndexList) length() int {
	return len(list.indexList)
}

func (list *IndexList) removeByItemId(id int) error {
	idx := list.indexOfId(id)
	if idx == -1 {
		return ErrLog("Couldn't remove item", nil)
	}

	ridx := list.indexList[idx].index

	list.indexList = append(list.indexList[:idx], list.indexList[idx+1:]...)

	for _, item := range list.indexList {
		if item.index > ridx {
			item.index--
		}
	}
	return nil
}

func (list *IndexList) indexOfId(id int) int {
	for i, item := range list.indexList {
		if item.value.Id == id {
			return i
		}
	}
	return -1
}

func (list *IndexList) editMessageById(id int, newContent string) error {
	idx := list.indexOfId(id)
	if idx == -1 {
		return ErrLog("Could not edit message using Id, id not found", nil)
	}

	list.indexList[idx].value.Content = newContent
	return nil
}

func (list *IndexList) moveToNewIndexById(id int, newIndex int) error {
	idx := list.indexOfId(id)
	if idx == -1 {
		return ErrLog("Id not found when moving", nil)
	}

	for i, item := range list.indexList {
		if item.index >= newIndex && i != idx {
			item.index++
		}
	}

	list.indexList[idx].index = newIndex
	return nil
}

func (list *IndexList) writeToFile() error {
	err := os.Truncate(ShoppingFile, 0)
	if err != nil {
		return ErrLog("Failed to truncate shopping file", err)
	}

	file, err := os.OpenFile(
		ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return ErrLog("Failed to open shopping file for writing", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range list.indexList {
		err := writer.Write(item.Serialize())
		if err != nil {
			return ErrLog("Something went wrong writing to the file", err)
		}
	}

	writer.Flush()
	err = writer.Error()
	if err != nil {
		return ErrLog("The writer experienced an error when writing", err)
	}
	return nil
}

func NewIndexList() IndexList {
	return IndexList {
		make([]IndexNode, 0),
	}
}

func ReadFromFile() (IndexList, error) {
	file, err := os.OpenFile(
		ShoppingFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return IndexList{}, ErrLog(
			"Something went wrong when opening the file for reading", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return IndexList{}, ErrLog("Reader failed reading all records", err)
	}

	list := NewIndexList()

	for _, record := range records {
		deserialized, err := deserialize(record)
		fmt.Println("deserialized:", deserialized);
		if err != nil {
			return IndexList{}, ErrLog("Failed to deserialize record", err)
		}

		list.indexList = append(list.indexList, deserialized)
	}

	return list, nil
}

func deserialize(serial []string) (IndexNode, error) {
	index, err := strconv.Atoi(serial[3])
	if err != nil {
		return IndexNode{}, ErrLog("Failed to convert index from file", err)
	}

	itemId, err := strconv.Atoi(serial[0])
	if err != nil {
		return IndexNode{}, ErrLog("Failed to convert item id from file", err)
	}

	return IndexNode {
		index: index,
		value: ShoppingItem{
			Id: itemId,
			Content: serial[1],
			Action: serial[2],
		},
	}, nil
}
