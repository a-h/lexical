package parse

import (
	"bytes"
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

	buf := bytes.NewBuffer([]byte{})
	for {
		current := pi.Index()
		ds := delimiter(pi)
		if ds.Success {
			rewind(pi, int(pi.Index()-current-1))
			return Success(name, string(buf.Bytes()), ds.Error)
		}
		r, err := pi.Advance()
		buf.WriteRune(r)
		if err != nil {
			if err == io.EOF {
				return Success(name, string(buf.Bytes()), err)
			}
			return Failure(name, err)
		}
	}
}
