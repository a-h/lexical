package parse

import (
	"io"
)

//TODO: Create a generic Until function.

// StringUntil captures runes until the delimiter or EOF is encountered and returns a string.
func StringUntil(delimiter Function) Function {
	return func(pi Input) Result {
		return stringUntil(pi, delimiter)
	}
}

func stringUntil(pi Input, delimiter Function) Result {
	name := "string until delimiter"

	runes := make([]rune, 0)
	for {
		current := pi.Index()
		ds := delimiter(pi)
		if ds.Success {
			rewind(pi, int(pi.Index()-current))
			return Success(name, string(runes), ds.Error)
		}
		r, err := pi.Advance()
		runes = append(runes, r)
		if err != nil {
			if err == io.EOF {
				return Success(name, string(runes), err)
			}
			return Failure(name, err)
		}
	}
}
