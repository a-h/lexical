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
			a:        func(pi Input) Result { return AnyRuneIn(pi, "A") },
			b:        func(pi Input) Result { return AnyRuneIn(pi, "a") },
			input:    "A",
			expected: true,
		},
		{
			name:     "'A' or 'a'",
			a:        func(pi Input) Result { return AnyRuneIn(pi, "A") },
			b:        func(pi Input) Result { return AnyRuneIn(pi, "a") },
			input:    "a",
			expected: true,
		},
		{
			name:     "'A' or 'a'",
			a:        func(pi Input) Result { return AnyRuneIn(pi, "A") },
			b:        func(pi Input) Result { return AnyRuneIn(pi, "a") },
			input:    "c",
			expected: false,
		},
	}

	for _, test := range tests {
		pi := input.NewFromString(test.name, test.input)
		result := Or(pi, test.a, test.b)
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
