package parse

import "fmt"

func Rune(pi Input, r rune) Result {
	name := fmt.Sprintf("rune '%v'", r)

	pr, err := pi.Peek()
	if pr == r {
		_, err = pi.Advance()
		return Success(name, pr, nil, err)
	}
	return Failure(name, fmt.Errorf("Expected '%v', but got '%v'", r, pr))
}
