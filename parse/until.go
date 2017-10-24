package parse

import (
	"bytes"
)

func StringUntil(pi Input, delimiter Function) Result {
	name := "string until delimiter"

	buf := bytes.NewBuffer([]byte{})
	for {
		ds := delimiter(pi)
		if ds.Success {
			return Success(name, string(buf.Bytes()), nil, ds.Error)
		}
		r, err := pi.Advance()
		if err != nil {
			return Failure(name, err)
		}
		buf.WriteRune(r)
	}
}

//TODO: More general version which accepts a parse function for capture and delimiter.
// It will need to separate the lexing and conversion into objects.

/*
func Until(pi Input, capture Function, until Function) Result {
	name := "until"

	buf := bytes.NewBuffer([]byte{})
	for {
		matchesCapture := capture(pi)
		if matchesCapture.Success {
			//TODO: Not possible to assume that we're going to get a string.
			// Maybe the capture and conversion need to be separated so that AnyRune / Or etc. only ever
			// capture strings and there's two levels, the parsing of parts of the structure, and converting
			// those captured strings into objects.
			buf.WriteString(string(matchesCapture.Item))
		}
		matchesUntil := until(pi)

		ds := delimiter(pi)
		if ds.Success {
			return Success(name, string(buf.Bytes()), nil, ds.Error)
		}
		r, err := pi.Advance()
		if err != nil {
			return Failure(name, err)
		}
		buf.WriteRune(r)
	}
}
*/
