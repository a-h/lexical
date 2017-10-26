package parse

// AnyRune returns a parser which will parse any rune at all.
func AnyRune() Function {
	return anyRune
}

func anyRune(pi Input) Result {
	r, err := pi.Advance()
	if err != nil {
		return Failure("any rune", err)
	}
	return Success("any rune", r, err)
}
