package parse

import (
	"io"
	"strings"
)

// StringUntil captures runes until the delimiter is encountered and returns a string.
func StringUntil(delimiter Function) Function {
	return func(pi Input) Result {
		return stringUntil(pi, delimiter, false)
	}
}

func StringUntilDelimiterOrEOF(delimiter Function) Function {
	return func(pi Input) Result {
		return stringUntil(pi, delimiter, true)
	}
}

func stringUntil(pi Input, delimiter Function, successOnEOF bool) Result {
	name := "string until delimiter"

	var sb strings.Builder
	for {
		current := pi.Index()
		ds := delimiter(pi)
		if ds.Success {
			rewind(pi, int(pi.Index()-current))
			return Success(name, sb.String(), ds.Error)
		}
		r, err := pi.Advance()
		if err != nil {
			if err == io.EOF && successOnEOF {
				return Success(name, sb.String(), nil)
			}
			return Failure(name, err)
		}
		sb.WriteRune(r)
	}
}
