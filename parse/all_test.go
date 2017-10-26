package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name             string
		parsers          []Function
		input            string
		expected         bool
		expectedPosition int64
	}{
		{
			name:             "A then B - success",
			parsers:          []Function{Rune('A'), Rune('B')},
			input:            "AB",
			expected:         true,
			expectedPosition: 2,
		},
		{
			name:             "A then B - failure",
			parsers:          []Function{Rune('A'), Rune('B')},
			input:            "AC",
			expected:         false,
			expectedPosition: 0,
		},
		{
			name:             "Any two runes",
			parsers:          []Function{AnyRune(), AnyRune()},
			input:            "ZC",
			expected:         true,
			expectedPosition: 2,
		},
	}

	for _, test := range tests {
		pi := input.NewFromString(test.name, test.input)
		parser := All(ConcatenateStringsCombiner, test.parsers...)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", test.name, test.input, test.expected, actual)
		}
		if test.expectedPosition != pi.Current {
			t.Errorf("test %v: for input '%v' expected to be at position %v but was at %v", test.name, test.input, test.expectedPosition, pi.Current)
		}
	}
}
