package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestStringUntil(t *testing.T) {
	tests := []struct {
		input           string
		parser          Function
		expected        bool
		expectedCapture string
	}{
		{
			input:           "name=value",
			parser:          StringUntil(String("=")),
			expected:        true,
			expectedCapture: "name",
		},
		{
			input:           "name value",
			parser:          StringUntil(Rune(' ')),
			expected:        true,
			expectedCapture: "name",
		},
		{
			input:           "<tag>",
			parser:          StringUntil(Rune('>')),
			expected:        true,
			expectedCapture: "<tag",
		},
		{
			input:           "this is a test",
			parser:          StringUntil(Rune('>')),
			expected:        true,
			expectedCapture: "this is a test",
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		result := test.parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v, catpured '%v'", i, test.input, test.expected, actual, result.Item)
		}
	}
}
