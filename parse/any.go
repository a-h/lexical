package parse

// Any returns the first match out of the parse functions passed in, or a failure if no
// parsers match.
func Any(functions ...Function) Function {
	return func(pi Input) Result {
		return any(pi, functions...)
	}
}

// Or returns the first or a or b. It's equivalent to the Any function with two parameters.
func Or(a, b Function) Function {
	return Any(a, b)
}

func any(pi Input, functions ...Function) Result {
	for _, f := range functions {
		r := f(pi)
		if r.Success {
			return r
		}
	}
	return Failure("any", nil)
}
