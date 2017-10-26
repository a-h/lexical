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
type MultipleResultCombiner func([]interface{}) interface{}

// ConcatenateStringsCombiner is a MultipleResultCombiner which concatenates the results together as a string.
var ConcatenateStringsCombiner MultipleResultCombiner = func(inputs []interface{}) interface{} {
	buf := bytes.NewBuffer([]byte{})
	for _, ip := range inputs {
		buf.WriteString(fmt.Sprintf("%v", ip))
	}
	return buf.String()
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
	return Success("inorder", combiner(results), nil)
}
