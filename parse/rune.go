package parse

import "fmt"

// Rune captures a single, specified rune.
func Rune(r rune) Function {
	return func(pi Input) Result {
		return parseRune(pi, r)
	}
}

func parseRune(pi Input, r rune) Result {
	name := "rune '" + string(r) + "'"

	pr, err := pi.Peek()
	if pr == r {
		_, err = pi.Advance()
		return Success(name, pr, err)
	}
	return Failure(name, fmt.Errorf("Expected '%v', but got '%v' with error %v", string(r), string(pr), err))
}
