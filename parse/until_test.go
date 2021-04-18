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
			expected:        false,
			expectedCapture: "",
		},
		{
			input:           "this is> a test",
			parser:          StringUntil(Rune('>')),
			expected:        true,
			expectedCapture: "this is",
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

func TestStringUntilInput(t *testing.T) {
	pi := input.NewFromString("abcd")
	result := StringUntil(String("cd"))(pi)
	if !result.Success {
		t.Errorf("unexpected failure")
	}
	if pi.Index() != 2 {
		t.Errorf("expected to be at index 2, but at %v", pi.Index())
	}
}

func BenchmarkUntil(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := StringUntil(Rune('Z'))
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
