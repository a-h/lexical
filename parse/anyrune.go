package parse

func AnyRune(pi Input) Result {
	r, err := pi.Advance()
	return Success("any rune", r, nil, err)
}
