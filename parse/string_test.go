package parse

import (
	"fmt"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestString(t *testing.T) {
	tests := []struct {
		input           string
		s               string
		expected        bool
		expectedCapture string
	}{
		{
			input:           "ABC",
			s:               "ABC",
			expected:        true,
			expectedCapture: "ABC",
		},
		{
			input:    "ABC",
			s:        "DEF",
			expected: false,
		},
		{
			input:    "ABC",
			s:        "BCD",
			expected: false,
		},
		{
			input:    "ABD",
			s:        "ABC",
			expected: false,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		parser := String(test.s)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		var expectedPosition int64
		if test.expected {
			expectedPosition = int64(len(test.input))
		}
		if pi.Current != expectedPosition {
			t.Errorf("test %v: for input '%v' expected to be at position %v but was at %v", i, test.input, expectedPosition, pi.Current)
		}
	}
}
