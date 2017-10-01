package lexical

import "testing"

func TestItemEquality(t *testing.T) {
	tests := []struct {
		a        Item
		b        Item
		expected bool
	}{
		{
			a:        Item{Type: ItemTypeError, Value: "123"},
			b:        Item{Type: ItemTypeError, Value: "123"},
			expected: true,
		},
		{
			a:        Item{Type: ItemTypeEOF, Value: "123"},
			b:        Item{Type: ItemTypeError, Value: "123"},
			expected: false,
		},
		{
			a:        Item{Type: ItemTypeEOF, Value: "123"},
			b:        Item{Type: ItemTypeEOF, Value: ""},
			expected: false,
		},
	}

	for i, test := range tests {
		actual := test.a.Eq(test.b)
		if actual != test.expected {
			t.Errorf("test %v: expected %v but got %v", i, test.expected, actual)
		}
	}
}

func TestItemString(t *testing.T) {
	tests := []struct {
		input    Item
		expected string
	}{
		{
			input:    Item{Type: ItemTypeError, Value: "Error message"},
			expected: "Error message",
		},
		{
			input:    Item{Type: ItemTypeEOF, Value: "Doesn't matter"},
			expected: "EOF",
		},
		{
			input:    Item{Type: ItemType(666), Value: "Hello"},
			expected: "Hello",
		},
		{
			input:    Item{Type: ItemType(666), Value: "Don't forget your lucky number."},
			expected: "Don't forg...",
		},
	}

	for i, test := range tests {
		actual := test.input.String()
		if actual != test.expected {
			t.Errorf("test %v: expected '%v', but got '%v'", i, test.expected, actual)
		}
	}
}
