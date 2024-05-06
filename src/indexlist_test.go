package main

import (
	"math/rand/v2"
	"reflect"
	"testing"
)

func NewShoppingItem(id int, content string, action string) ShoppingItem {
	return ShoppingItem {
		Id: id,
		Content: content,
		Action: action,
	}
}

const letters = "abcdefghijklmnopqrstuwvxzy"
const randLength = 10

func MakeRandomStringLen(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

func CreateRandomShoppingItems(count int) []ShoppingItem {
	itemList := make([]ShoppingItem, count)

	for i := 0; i < count; i++ {
		itemList[i] = NewShoppingItem(rand.IntN(100), MakeRandomStringLen(randLength), MakeRandomStringLen(randLength))
	}

	return itemList
}

func NewFilledNode(index int, id int, content string, action string) IndexNode {
	return IndexNode {
		index: index,
		value: NewShoppingItem(id, content, action),
	}
}

/*
 * Serialisation tests
 */
func TestSerialisationFilled(t *testing.T) {
	node := NewFilledNode(0, 0, "main", "add")
	want := []string{"0", "main", "add", "0"}

	if !reflect.DeepEqual(node.Serialize(), want) {
		t.Fatal("Serialize not equal to expected outcome")
	}
}


/*
 * Deserialisation Tests
 */
func TestSerialiseThenDeserialize(t *testing.T) {
	node := NewFilledNode(0, 0, "something", "add")

	deser, err := Deserialize(node.Serialize())

	if err != nil || deser != node {
		t.Fatal("Serializsing and deserialising failed")
	}
}

/*
 * Length Tests
 */
func TestIndexListLen0(t *testing.T) {
	list := NewIndexList()
	if list.Length() != 0 {
		t.Fatal("Length was not 0")
	}
}

func TestIndexListLen5(t *testing.T) {
	list := NewIndexList()
	list.indexList = make([]IndexNode, 5)

	if list.Length() != 5 {
		t.Fatal("Length was not 5")
	}
}

/*
 * IndexList add tests
 */
func TestIndexListAddSingle(t *testing.T) {
	list := NewIndexList()
	tnode := NewShoppingItem(0, "first", "add")
	list.add(tnode)

	if list.indexList[0].index != 0 {
		t.Fatal("Wrong index on first element")
	}

	if list.indexList[0].value != tnode {
		t.Fatal("Bad first node added")
	}
}

func TestIndexListAddThree(t *testing.T) {
	list := NewIndexList()
	node1 := NewShoppingItem(0, "first", "add")
	node2 := NewShoppingItem(1, "second", "add")
	node3 := NewShoppingItem(2, "third", "add")
	list.add(node1)
	list.add(node2)
	list.add(node3)

	for i, item := range list.indexList {
		if item.index != i {
			t.Fatal("Wrong index on indexlist item after adding")
		}
	}

	if list.indexList[0].value != node1 {
		t.Fatal("Node does not match added")
	}
	if list.indexList[1].value != node2 {
		t.Fatal("Node does not match added")
	}
	if list.indexList[2].value != node3 {
		t.Fatal("Node does not match added")
	}
}

/* 
 * RemoveItemById test
 */
func TestRemoveSingleItem(t *testing.T) {
	items := CreateRandomShoppingItems(10)
	list := NewIndexList()

	for _, item := range items {
		list.add(item)
	}

	// the string generators cannot make this, so I will remove this one
	node := NewShoppingItem(16, "--", "--")
	list.add(node)
	list.RemoveByItemId(16)

	for _, item := range list.indexList {
		if item.value == node {
			t.Fatal("Failed to remove by id")
		}
	}

}
