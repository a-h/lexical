package parse

import (
	"bytes"
	"fmt"
)

// All ensures that all of the parsers are captured, or winds the whole set of captures back.
func All(combiner MultipleResultCombiner, functions ...Function) Function {
	return func(pi Input) Result {
		return all(pi, combiner, functions...)
	}
}

// MultipleResultCombiner combines the results from multiple parse operations into a single result.
type MultipleResultCombiner func([]interface{}) (result interface{}, ok bool)

// ConcatenateToStringCombiner is a MultipleResultCombiner which concatenates the results together as a string.
var ConcatenateToStringCombiner MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
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

func all(pi Input, combiner MultipleResultCombiner, functions ...Function) Result {
	results := make([]interface{}, len(functions))
	start := pi.Index()
	for i, f := range functions {
		r := f(pi)
		if !r.Success {
			rewind(pi, int(pi.Index()-start))
			return r
		}
		results[i] = r.Item
	}

	// Combine all the results using the provided function.
	item, ok := combiner(results)
	if !ok {
		return Failure("inorder", fmt.Errorf("failed to combine results"))
	}
	return Success("inorder", item, nil)
}
