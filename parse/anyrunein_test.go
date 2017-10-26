package parse

import (
	"fmt"
	"io"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestAnyRuneIn(t *testing.T) {
	tests := []struct {
		input         string
		in            string
		expected      bool
		expectedError error
	}{
		{
			input:    "A",
			in:       "ABC",
			expected: true,
		},
		{
			input:    "B",
			in:       "ABC",
			expected: true,
		},
		{
			input:    "C",
			in:       "ABC",
			expected: true,
		},
		{
			input:    "D",
			in:       "ABC",
			expected: false,
		},
		{
			input:         "",
			in:            "A",
			expected:      false,
			expectedError: io.EOF,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		result := AnyRuneIn(pi, test.in)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if result.Error != test.expectedError {
			t.Errorf("test %v: for input '%v' expected error '%v' but got '%v'", i, test.input, test.expectedError, result.Error)
		}
	}
}