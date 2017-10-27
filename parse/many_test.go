package parse

import (
	"fmt"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestMany(t *testing.T) {
	tests := []struct {
		input         string
		parser        Function
		expectedMatch bool
		expectedItem  string
	}{
		{
			input:         "AAAAAAAA",
			parser:        Many(Rune('A'), ConcatenateToStringCombiner, 0, 500),
			expectedMatch: true,
			expectedItem:  "AAAAAAAA",
		},
		{
			input:         "AAABBB",
			parser:        Many(Rune('A'), ConcatenateToStringCombiner, 0, 500),
			expectedMatch: true,
			expectedItem:  "AAA",
		},
		{
			input:         "AAABBB",
			parser:        Many(Or(Rune('A'), Rune('B')), ConcatenateToStringCombiner, 0, 500),
			expectedMatch: true,
			expectedItem:  "AAABBB",
		},
		{
			input:         "AAABBB",
			parser:        Many(Rune('A'), ConcatenateToStringCombiner, 4, 500),
			expectedMatch: false,
		},
		{
			input:         "AAABBB",
			parser:        Many(Rune('A'), ConcatenateToStringCombiner, 1, 2),
			expectedMatch: true,
			expectedItem:  "AA",
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		result := test.parser(pi)
		actualMatch := result.Success
		if actualMatch != test.expectedMatch {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expectedMatch, actualMatch)
		}
		if test.expectedMatch && result.Item != test.expectedItem {
			t.Errorf("test %v: for input '%v' expected to capture '%v' but got '%v'", i, test.input, test.expectedItem, result.Item)
		}
	}
}
