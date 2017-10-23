package parse

import (
	"fmt"
)

func String(pi Input, s string) Result {
	name := fmt.Sprintf("string: '%v'", s)

	advancedBy := 0
	for _, sr := range s {
		pr, err := pi.Peek()
		if pr != sr {
			err = rewind(pi, advancedBy)
			return Failure(name, err)
		}
		pi.Advance()
		advancedBy++
	}
	return Success(name, s, nil, nil)
}

func rewind(pi Input, times int) (err error) {
	for i := 0; i < times; i++ {
		_, err = pi.Retreat()
		if err != nil {
			return
		}
	}
	return
}
