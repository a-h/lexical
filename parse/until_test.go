package parse

import (
	"fmt"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestStringUntil(t *testing.T) {
	tests := []struct {
		input           string
		s               string
		expected        bool
		expectedCapture string
	}{
		{
			input:           "name=value",
			s:               "=",
			expected:        true,
			expectedCapture: "name",
		},
		{
			input:           "name value",
			s:               " ",
			expected:        true,
			expectedCapture: "name",
		},
		{
			input:           "<tag>",
			s:               ">",
			expected:        true,
			expectedCapture: "<tag",
		},
		{
			input:           "this is a test",
			s:               ">",
			expected:        false,
			expectedCapture: "this is a test",
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		result := StringUntil(pi, func(ppi Input) Result { return String(ppi, test.s) })
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
	}
}
