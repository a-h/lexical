package parse

import "errors"

// Then executes one function, then another, comining the results using the provided function.
func Then(a, b Function, combiner MultipleResultCombiner) Function {
	return func(pi Input) Result {
		return then(pi, a, b, combiner)
	}
}

func then(pi Input, a, b Function, combiner MultipleResultCombiner) Result {
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
