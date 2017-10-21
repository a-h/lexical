package parse

import "fmt"

type InterfaceResultItem struct {
	name  string
	value interface{}
}

func (item InterfaceResultItem) Name() string {
	return item.name
}

func (item InterfaceResultItem) Value() interface{} {
	return item.value
}

func (item InterfaceResultItem) Eq(cmp ResultItem) bool {
	if cmp.Name() != item.Name() {
		return false
	}
	if cmp.Value() != item.Value() {
		return false
	}
	return true
}

// String returns the string representation of an item, truncated to 10 characters.
func (item InterfaceResultItem) String() string {
	v := fmt.Sprintf("%v", item.value)
	if len(v) > 13 {
		return fmt.Sprintf("%v: %v...", item.name, v[0:10])
	}
	return fmt.Sprintf("%v: %v", item.name, v)
}
