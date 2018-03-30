package parse

// String captures a specific string.
func String(s string) Function {
	return func(pi Input) Result {
		return parseString(pi, s)
	}
}

func parseString(pi Input, s string) Result {
	name := "string: '" + s + "'"

	advancedBy := 0
	for _, sr := range s {
		pr, err := pi.Peek()
		if pr != sr {
			rewind(pi, advancedBy)
			return Failure(name, err)
		}
		pi.Advance()
		advancedBy++
	}
	return Success(name, s, nil)
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
