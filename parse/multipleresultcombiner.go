package parse

import (
	"fmt"
	"strconv"
)

// MultipleResultCombiner combines the results from multiple parse operations into a single result.
type MultipleResultCombiner func([]interface{}) (result interface{}, ok bool)

// WithStringConcatCombiner is a MultipleResultCombiner which concatenates the results together as a string.
func WithStringConcatCombiner(inputs []interface{}) (interface{}, bool) {
	var buf []byte
	for _, ip := range inputs {
		switch v := ip.(type) {
		case rune:
			buf = append(buf, string(v)...)
		case string:
			buf = append(buf, v...)
		default:
			buf = append(buf, fmt.Sprintf("%v", v)...)
		}
	}
	return string(buf), true
}

// WithIntegerCombiner is a MultipleResultCombiner which concatenates the results together as a string then parses
// the result as an integer.
func WithIntegerCombiner(items []interface{}) (item interface{}, value bool) {
	si, ok := WithStringConcatCombiner(items)
	if !ok {
		return 0, false
	}
	s, _ := si.(string)
	i, err := strconv.ParseInt(s, 10, 32)
	return int(i), err == nil
}
