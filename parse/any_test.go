package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestOr(t *testing.T) {
	tests := []struct {
		name     string
		a        Function
		b        Function
		input    string
		expected bool
	}{
		{
			name:     "'A' or 'a'",
			a:        Rune('A'),
			b:        Rune('a'),
			input:    "A",
			expected: true,
		},
		{
			name:     "'A' or 'a'",
			a:        Rune('A'),
			b:        Rune('a'),
			input:    "a",
			expected: true,
		},
		{
			name:     "'A' or 'a'",
			a:        Rune('A'),
			b:        Rune('a'),
			input:    "c",
			expected: false,
		},
	}

	for _, test := range tests {
		pi := input.NewFromString(test.name, test.input)
		parser := Or(test.a, test.b)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", test.name, test.input, test.expected, actual)
		}
		var expectedPosition int64
		if test.expected {
			expectedPosition = 1
		}
		if pi.Current != expectedPosition {
			t.Errorf("test %v: for input '%v' expected to be at position %v but was at %v", test.name, test.input, expectedPosition, pi.Current)
		}
	}
}
