package lexical

import (
	"fmt"
)

// An Item represents a Lexical Item consisting of its type and its value.
type Item struct {
	// The Type of the item, e.g. "Number", or "LessThan".
	Type ItemType
	// Value is the string value of the Item, e.g. "1234" or "<".
	Value string
}

// Eq compares two items and returns true if they are equal.
func (item Item) Eq(cmp Item) bool {
	return item.Type == cmp.Type && item.Value == cmp.Value
}

// String returns the string representation of an item, truncated to 10 characters.
func (item Item) String() string {
	switch item.Type {
	case ItemTypeEOF:
		return "EOF"
	case ItemTypeError:
		return fmt.Sprintf("err: %v", item.Value)
	}
	if len(item.Value) > 13 {
		return fmt.Sprintf("%v...", item.Value[0:10])
	}
	return string(item.Value)
}

// An ItemType is the type of an Item, e.g. "Number", or "LessThan".
type ItemType int

const (
	// ItemTypeError defines an error state.
	ItemTypeError ItemType = iota
	// ItemTypeEOF defines the end of the file.
	ItemTypeEOF
)
