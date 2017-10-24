package parse

import "testing"
import "io"

func TestResultEq(t *testing.T) {
	tests := []struct {
		a        Result
		b        Result
		expected bool
	}{
		{
			a:        Result{Name: "a", Item: "123"},
			b:        Result{Name: "a", Item: "123"},
			expected: true,
		},
		{
			a:        Result{Name: "a", Item: "123"},
			b:        Result{Name: "b", Item: "123"},
			expected: false,
		},
		{
			a:        Result{Name: "a", Item: "123"},
			b:        Result{Name: "a", Item: ""},
			expected: false,
		},
		{
			a:        Result{Name: "a", Item: "123"},
			b:        Result{Name: "a", Item: 123},
			expected: false,
		},
	}

	for i, test := range tests {
		actual := test.a.Eq(test.b)
		if actual != test.expected {
			t.Errorf("test %v: expected %v but got %v", i, test.expected, actual)
		}
	}
}

func TestResultString(t *testing.T) {
	tests := []struct {
		input    Result
		expected string
	}{
		{
			input:    Success("a", 123, nil, nil),
			expected: "✓ (a) 123",
		},
		{
			input:    Failure("a", nil),
			expected: "✗ (a) err: <nil>",
		},
		{
			input:    Success("a", "Don't forget your lucky number.", nil, nil),
			expected: "✓ (a) Don't forg...",
		},
		{
			input:    Failure("a", io.EOF),
			expected: "✗ (a) err: EOF",
		},
		{
			input:    Success("a", 123, nil, io.EOF),
			expected: "✓ (a) 123\n✓ (a) err: EOF",
		},
	}

	for i, test := range tests {
		actual := test.input.String()
		if actual != test.expected {
			t.Errorf("test %v: expected '%v', but got '%v'", i, test.expected, actual)
		}
	}
}
