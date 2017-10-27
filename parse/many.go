package parse

import "errors"

// Many captures the function at least x times and at most y times and sets the
// result item to an array of the function captures.
func Many(combiner MultipleResultCombiner, atLeast, atMost int, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, atLeast, atMost, f)
	}
}

func many(pi Input, combiner MultipleResultCombiner, atLeast, atMost int, f Function) Result {
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
