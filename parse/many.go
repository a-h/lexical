package parse

import "errors"

// Many captures the function at least x times and at most y times and sets the
// result item to an array of the function captures.
func Many(combiner MultipleResultCombiner, atLeast, atMost int, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, atLeast, atMost, f)
	}
}

// Times captures the parser function a set number of times.
func Times(combiner MultipleResultCombiner, times int, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, times, times, f)
	}
}

// AtLeast captures the passed function at least the number of times provided.
func AtLeast(combiner MultipleResultCombiner, times int, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, times, -1, f)
	}
}

// AtMost captures the passed function between one and the number of times provided.
func AtMost(combiner MultipleResultCombiner, times int, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, 1, times, f)
	}
}

// Optional provides an optional parser.
func Optional(combiner MultipleResultCombiner, f Function) Function {
	return func(pi Input) Result {
		return many(pi, combiner, 0, 1, f)
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
		if atMost > 0 && len(results) == atMost {
			break
		}
	}

	item, ok := combiner(results)
	if !ok {
		return Failure("many", errors.New("failed to combine results"))
	}
	return Success("many", item, nil)
}
