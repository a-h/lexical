package parse

// Then executes one function, then another, comining the results using the provided function.
func Then(a, b Function, mapper ResultCombiner) Function {
	return func(pi Input) Result {
		return then(pi, a, b, mapper)
	}
}

// ResultCombiner merges the results from a and b
type ResultCombiner func(a, b interface{}) interface{}

func then(pi Input, a, b Function, mapper ResultCombiner) Result {
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

	return Success("then", mapper(ar.Item, br.Item), br.Error)
}
