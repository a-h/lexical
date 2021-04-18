package parse

import (
	"fmt"
)

// All ensures that all of the parsers are captured, or winds the whole set of captures back.
func All(combiner MultipleResultCombiner, functions ...Function) Function {
	return func(pi Input) Result {
		return all(pi, combiner, functions...)
	}
}

func all(pi Input, combiner MultipleResultCombiner, functions ...Function) Result {
	results := make([]interface{}, len(functions))
	start := pi.Index()
	for i := 0; i < len(functions); i++ {
		r := functions[i](pi)
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
