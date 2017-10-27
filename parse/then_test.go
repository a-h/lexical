package parse

import (
	"fmt"
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
			parser:          Then(Rune('A'), Rune('B'), ConcatenateToStringCombiner),
			expected:        true,
			expectedCapture: "AB",
		},
		{
			input:    "ab",
			parser:   Then(Rune('A'), Rune('B'), ConcatenateToStringCombiner),
			expected: false,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		result := test.parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v, catpured '%v'", i, test.input, test.expected, actual, result.Item)
		}
	}
}
