package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestThen(t *testing.T) {
	tests := []struct {
		input           string
		parser          Function
		expected        bool
		expectedCapture string
	}{
		{
			input:           "AB",
			parser:          Then(WithStringConcatCombiner, Rune('A'), Rune('B')),
			expected:        true,
			expectedCapture: "AB",
		},
		{
			input:    "ab",
			parser:   Then(WithStringConcatCombiner, Rune('A'), Rune('B')),
			expected: false,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		result := test.parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v, catpured '%v'", i, test.input, test.expected, actual, result.Item)
		}
	}
}
