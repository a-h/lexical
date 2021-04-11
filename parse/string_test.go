package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestString(t *testing.T) {
	tests := []struct {
		input           string
		s               string
		expected        bool
		expectedCapture string
		expectedIndex   int64
	}{
		{
			input:           "ABC",
			s:               "ABC",
			expected:        true,
			expectedCapture: "ABC",
			expectedIndex:   3,
		},
		{
			input:           "ABC",
			s:               "AB",
			expected:        true,
			expectedCapture: "AB",
			expectedIndex:   2,
		},
		{
			input:         "ABC",
			s:             "DEF",
			expected:      false,
			expectedIndex: 0,
		},
		{
			input:         "ABC",
			s:             "BCD",
			expected:      false,
			expectedIndex: 0,
		},
		{
			input:         "ABD",
			s:             "ABC",
			expected:      false,
			expectedIndex: 0,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		parser := String(test.s)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if test.expectedIndex != pi.Index() {
			t.Errorf("test %v: for input '%v' expected index %d, got %d", i, test.input, test.expectedIndex, pi.Index())
		}
	}
}

func BenchmarkString(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		String("ABCDEFG")(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
