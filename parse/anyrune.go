package parse

func AnyRune(pi Input) Result {
	r, err := pi.Advance()
	if err != nil {
		return Failure("any rune", err)
	}
	return Success("any rune", r, nil, err)
}
