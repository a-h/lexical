package parse

import (
	"errors"
	"fmt"
)

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

	globalRollback := pi.Index()
	for {
		localRollback := pi.Index()
		r := f(pi)
		if !r.Success {
			rewind(pi, int(pi.Index()-localRollback))
			break
		}
		results = append(results, r.Item)
		if atMost > 0 && len(results) == atMost {
			break
		}
	}
	if len(results) < atLeast {
		rewind(pi, int(pi.Index()-globalRollback))
		return Failure("many", fmt.Errorf("expected at least %d results, got %d", atLeast, len(results)))
	}

	item, ok := combiner(results)
	if !ok {
		return Failure("many", errors.New("failed to combine results"))
	}
	return Success("many", item, nil)
}
