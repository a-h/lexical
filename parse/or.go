package parse

// Or returns the first match out of the parse functions passed in, or a failure if no
// parsers match.
func Or(pi Input, functions ...Function) Result {
	for _, f := range functions {
		r := f(pi)
		if r.Success {
			return r
		}
	}
	return Failure("or", nil)
}
