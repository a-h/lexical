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
	}

	for _, test := range tests {
		pi := input.NewFromString(test.name, test.input)
		parser := All(ConcatenateToStringCombiner, test.parsers...)
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

func TestConcatenateToStringCombiner(t *testing.T) {
	inputs := []interface{}{
		'A',
		"BCD",
		'E',
		1,
		2.0,
	}
	result, _ := ConcatenateToStringCombiner(inputs)
	if result != "ABCDE12" {
		t.Errorf("Expected 'ABCDE12.0', but got '%v'", result)
	}
}
