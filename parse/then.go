package parse

import (
	"errors"
)

// Then executes one function, then another, comining the results using the provided function.
func Then(combiner MultipleResultCombiner, a, b Function) Function {
	return func(pi Input) Result {
		return then(pi, combiner, a, b)
	}
}

func then(pi Input, combiner MultipleResultCombiner, a, b Function) Result {
	start := pi.Index()

	ar := a(pi)
	if !ar.Success {
		rewind(pi, int(pi.Index()-start))
		return ar
	}

	br := b(pi)
	if !br.Success {
		rewind(pi, int(pi.Index()-start))
		return br
	}

	item, ok := combiner([]interface{}{ar.Item, br.Item})
	if !ok {
		return Failure("then", errors.New("failed to combine results"))
	}
	return Success("then", item, br.Error)
}
