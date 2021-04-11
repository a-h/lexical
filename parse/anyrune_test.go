package parse

import (
	"io"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestAnyRune(t *testing.T) {
	tests := []struct {
		input         string
		expected      bool
		expectedError error
		expectedIndex int64
	}{
		{
			input:         "A",
			expected:      true,
			expectedIndex: 1,
		},
		{
			input:         "",
			expected:      false,
			expectedError: io.EOF,
			expectedIndex: 1,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		parser := AnyRune()
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if result.Error != test.expectedError {
			t.Errorf("test %v: for input '%v' expected error '%v' but got '%v'", i, test.input, test.expectedError, result.Error)
		}
		if test.expectedIndex != pi.Index() {
			t.Errorf("test %v: for input '%v' expected index %d, got %d", i, test.input, test.expectedIndex, pi.Index())
		}
	}
}

func BenchmarkAnyRune(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := AnyRune()
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
