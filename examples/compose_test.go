package examples

import (
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestCompose(t *testing.T) {
	parser := parse.Or(parse.Rune('A'), parse.Rune('B'))

	matchesA := parser(input.NewFromString("A")).Success // true
	matchesB := parser(input.NewFromString("B")).Success // true
	matchesC := parser(input.NewFromString("C")).Success // false

	if !matchesA {
		t.Errorf("for 'A', expected true, got false")
	}
	if !matchesB {
		t.Errorf("for 'B', expected true, got false")
	}
	if matchesC {
		t.Errorf("for 'C', expected false, got true")
	}
}

func TestMany(t *testing.T) {
	// parse.WithIntegerCombiner concatentates the captured runes into a string,
	// and parses the result to an integer.
	oneToThreeNumbers := parse.Many(parse.WithIntegerCombiner,
		1, // minimum match count
		3, // maximum match count
		parse.ZeroToNine)

	resultA := oneToThreeNumbers(input.NewFromString("123"))
	if !resultA.Success {
		t.Error("for '123' expected success to be true, got false")
	}
	if resultA.Item != 123 {
		t.Errorf("for '123' expected value of 123, but got '%v'", resultA.Item)
	}

	resultB := oneToThreeNumbers(input.NewFromString("1234"))
	if !resultB.Success {
		t.Errorf("for '1234', expected success to be true, got false")
	}
	if resultB.Item != 123 {
		t.Errorf("for '1234' expected value of 123, but got '%v'", resultA.Item)
	}

	// This Many function will stop reading at the 'a'.
	resultC := oneToThreeNumbers(input.NewFromString("1a234"))
	if !resultC.Success {
		t.Errorf("for '1a234', expected success to be true, got false")
	}
	if resultC.Item != 1 {
		t.Errorf("for '1a234' expected value of 1, but got '%v'", resultA.Item)
	}

	// Parse letters into a string
	upToThreeLetters := parse.AtMost(parse.WithStringConcatCombiner, 3, parse.Letter)
	letters := upToThreeLetters(input.NewFromString("ABC1"))
	resultItem, ok := letters.Item.(string)
	if !ok || resultItem != "ABC" {
		t.Errorf("for 'ABC1', expected to extract 'ABC', but extracted '%v'", letters.Item)
	}
}
