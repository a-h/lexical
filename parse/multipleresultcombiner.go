package parse

import (
	"bytes"
	"fmt"
	"strconv"
)

// MultipleResultCombiner combines the results from multiple parse operations into a single result.
type MultipleResultCombiner func([]interface{}) (result interface{}, ok bool)

// WithStringConcatCombiner is a MultipleResultCombiner which concatenates the results together as a string.
var WithStringConcatCombiner MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	buf := bytes.NewBuffer([]byte{})
	for _, ip := range inputs {
		switch v := ip.(type) {
		case rune:
			buf.WriteRune(v)
		case string:
			buf.WriteString(v)
		case Function:
			buf.WriteString("error: function passed to combiner")
			return buf.String(), false
		default:
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return buf.String(), true
}

// WithIntegerCombiner is a MultipleResultCombiner which concatenates the results together as a string then parses
// the result as an integer.
func WithIntegerCombiner(items []interface{}) (item interface{}, value bool) {
	s := ""
	for _, r := range items {
		switch v := r.(type) {
		case rune:
			s += string(v)
		case string:
			s += v
		default:
			return 0, false
		}
	}
	i, err := strconv.ParseInt(s, 10, 32)
	return int(i), err == nil
}
