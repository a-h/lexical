package parse

import (
	"fmt"
	"strings"
)

func AnyRuneIn(pi Input, in string) Result {
	name := fmt.Sprintf("any rune in '%v'", in)

	pr, err := pi.Peek()
	if strings.ContainsRune(in, pr) {
		_, err = pi.Advance()
		return Success(name, pr, nil, err)
	}
	return Failure(name, err)
}
