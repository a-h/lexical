package parse

import "errors"

// Many captures the function at least x times and at most y times and sets the
// result item to an array of the function captures.
func Many(f Function, combiner MultipleResultCombiner, atLeast, atMost int) Function {
	return func(pi Input) Result {
		return many(pi, f, combiner, atLeast, atMost)
	}
}

func many(pi Input, f Function, combiner MultipleResultCombiner, atLeast, atMost int) Result {
	results := make([]interface{}, 0)

	start := pi.Index()
	for {
		r := f(pi)
		if !r.Success {
			if len(results) < atLeast {
				// Roll back, because we didn't get enough.
				rewind(pi, int(pi.Index()-start))
				return r
			}
			// We're OK to stop.
			break
		}
		results = append(results, r.Item)
		if len(results) == atMost {
			break
		}
	}

	item, ok := combiner(results)
	if !ok {
		return Failure("many", errors.New("failed to combine results"))
	}
	return Success("many", item, nil)
}
