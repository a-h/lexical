package parse

import (
	"fmt"
	"strings"
	"unicode"
)

// RuneWhere captures a rune which matches a predicate.
func RuneWhere(predicate func(r rune) bool) Function {
	return func(pi Input) Result {
		return runeWhere(pi, "any rune where", predicate)
	}
}

// RuneIn captures a rune if it's within the input set.
func RuneIn(set string) Function {
	return func(pi Input) Result {
		name := fmt.Sprintf("any rune in '%v'", set)
		return runeWhere(pi, name, func(r rune) bool { return strings.ContainsRune(set, r) })
	}
}

// RuneNotIn captures a rune if it's not within the input set.
func RuneNotIn(set string) Function {
	return func(pi Input) Result {
		name := fmt.Sprintf("any rune not in '%v'", set)
		return runeWhere(pi, name, func(r rune) bool { return !strings.ContainsRune(set, r) })
	}
}

func runeWhere(pi Input, name string, predicate func(r rune) bool) Result {
	pr, err := pi.Peek()
	if predicate(pr) {
		_, err = pi.Advance()
		return Success(name, pr, err)
	}
	return Failure(name, err)
}

// RuneInRanges returns a parser which accepts a rune within the specified Unicode range.
func RuneInRanges(rts ...*unicode.RangeTable) Function {
	return func(pi Input) Result {
		return runeWhere(pi, "rune in ranges", func(r rune) bool { return unicode.IsOneOf(rts, r) })
	}
}

// Letter returns a parser which accepts a rune within the Letter Unicode range.
var Letter = RuneInRanges(unicode.Letter)

// ZeroToNine returns a parser which accepts a rune within range, i.e. 0-9.
var ZeroToNine = RuneIn("0123456789")
