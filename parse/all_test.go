package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name             string
		parsers          []Function
		input            string
		expected         bool
		expectedPosition int64
		expectedItem     string
	}{
		{
			name:             "A then B - success",
			parsers:          []Function{Rune('A'), Rune('B')},
			input:            "AB",
			expected:         true,
			expectedPosition: 2,
			expectedItem:     "AB",
		},
		{
			name:             "A then B - failure",
			parsers:          []Function{Rune('A'), Rune('B')},
			input:            "AC",
			expected:         false,
			expectedPosition: 0,
		},
		{
			name:             "Any two runes",
			parsers:          []Function{AnyRune(), AnyRune()},
			input:            "ZC",
			expected:         true,
			expectedPosition: 2,
			expectedItem:     "ZC",
		},
		{
			name:             "<abc>",
			parsers:          []Function{Rune('<'), String("abc"), Rune('>')},
			input:            "<abc>",
			expected:         true,
			expectedPosition: 5,
			expectedItem:     "<abc>",
		},
		{
			name:             "<abc> - until c",
			parsers:          []Function{Rune('<'), StringUntil(String("c")), String("c"), Rune('>')},
			input:            "<abc>",
			expected:         true,
			expectedPosition: 5,
			expectedItem:     "<abc>",
		},
	}

	for _, test := range tests {
		pi := input.NewFromString(test.input)
		parser := All(WithStringConcatCombiner, test.parsers...)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", test.name, test.input, test.expected, actual)
		}
		if test.expectedPosition != pi.Current {
			t.Errorf("test %v: for input '%v' expected to be at position %v but was at %v", test.name, test.input, test.expectedPosition, pi.Current)
		}
		if test.expected && test.expectedItem != result.Item {
			t.Errorf("test %v: for input '%v' expected item '%v' but was at %v", test.name, test.input, test.expectedItem, result.Item)
		}
	}
}

func BenchmarkAll(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := All(WithStringConcatCombiner, Rune('A'), Rune('B'))
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
