package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestThen(t *testing.T) {
	tests := []struct {
		input           string
		parser          Function
		expected        bool
		expectedCapture string
		expectedIndex   int64
	}{
		{
			input:           "AB",
			parser:          Then(WithStringConcatCombiner, Rune('A'), Rune('B')),
			expected:        true,
			expectedCapture: "AB",
			expectedIndex:   2,
		},
		{
			input:         "ab",
			parser:        Then(WithStringConcatCombiner, Rune('A'), Rune('B')),
			expected:      false,
			expectedIndex: 0,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		result := test.parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v, catpured '%v'", i, test.input, test.expected, actual, result.Item)
		}
		if test.expectedIndex != pi.Index() {
			t.Errorf("test %v: for input '%v' expected index %d, got %d", i, test.input, test.expectedIndex, pi.Index())
		}
	}
}

func BenchmarkThen(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := Then(WithStringConcatCombiner, Rune('A'), Rune('B'))
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
