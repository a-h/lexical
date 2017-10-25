package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestInOrder(t *testing.T) {
	tests := []struct {
		name             string
		parsers          []Function
		input            string
		expected         bool
		expectedPosition int64
	}{
		{
			name: "A then B - success",
			parsers: []Function{
				func(pi Input) Result { return Rune(pi, 'A') },
				func(pi Input) Result { return Rune(pi, 'B') },
			},
			input:            "AB",
			expected:         true,
			expectedPosition: 1,
		},
		{
			name: "A then B - failure",
			parsers: []Function{
				func(pi Input) Result { return Rune(pi, 'A') },
				func(pi Input) Result { return Rune(pi, 'B') },
			},
			input:            "AC",
			expected:         false,
			expectedPosition: -1,
		},
		{
			name: "Any two runes",
			parsers: []Function{
				func(pi Input) Result { return AnyRune(pi) },
				func(pi Input) Result { return AnyRune(pi) },
			},
			input:            "ZC",
			expected:         true,
			expectedPosition: 1,
		},
	}

	for _, test := range tests {
		pi := input.NewFromString(test.name, test.input)
		result := InOrder(pi, test.parsers...)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", test.name, test.input, test.expected, actual)
		}
		if test.expectedPosition != pi.Current {
			t.Errorf("test %v: for input '%v' expected to be at position %v but was at %v", test.name, test.input, test.expectedPosition, pi.Current)
		}
	}
}
