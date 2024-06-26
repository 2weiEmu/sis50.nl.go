package node

import (
	"strconv"

	s "github.com/2weiEmu/sis50.nl.go/pkg/shopping"
)

type IndexNode struct {
	Index int;
	Value s.ShoppingItem;
}
func (node *IndexNode) Serialize() []string {
	return append(node.Value.Serialize(), strconv.Itoa(node.Index))
}

