package parse

import (
	"io"
	"strings"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestRuneIn(t *testing.T) {
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
		pi := input.NewFromString(test.input)
		parser := RuneIn(test.in)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if result.Error != test.expectedError {
			t.Errorf("test %v: for input '%v' expected error '%v' but got '%v'", i, test.input, test.expectedError, result.Error)
		}
	}
}

func BenchmarkRuneWhere(b *testing.B) {
	isAToZ := func(r rune) bool { return strings.IndexRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", r) >= 0 }
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := RuneWhere(isAToZ)
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}

func BenchmarkRuneIn(b *testing.B) {
	b.ReportAllocs()
	parser := RuneIn("ABC")
	for n := 0; n < b.N; n++ {
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
