package parse

import "testing"

func TestInterfaceResultItemEq(t *testing.T) {
	tests := []struct {
		a        InterfaceResultItem
		b        InterfaceResultItem
		expected bool
	}{
		{
			a:        InterfaceResultItem{name: "a", value: "123"},
			b:        InterfaceResultItem{name: "a", value: "123"},
			expected: true,
		},
		{
			a:        InterfaceResultItem{name: "a", value: "123"},
			b:        InterfaceResultItem{name: "b", value: "123"},
			expected: false,
		},
		{
			a:        InterfaceResultItem{name: "a", value: "123"},
			b:        InterfaceResultItem{name: "a", value: ""},
			expected: false,
		},
		{
			a:        InterfaceResultItem{name: "a", value: "123"},
			b:        InterfaceResultItem{name: "a", value: 123},
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

func TestInterfaceResultItemString(t *testing.T) {
	tests := []struct {
		input    InterfaceResultItem
		expected string
	}{
		{
			input:    InterfaceResultItem{name: "a", value: "123"},
			expected: "a: 123",
		},
		{
			input:    InterfaceResultItem{name: "b", value: 123},
			expected: "b: 123",
		},
		{
			input:    InterfaceResultItem{name: "a", value: "Don't forget your lucky number."},
			expected: "a: Don't forg...",
		},
	}

	for i, test := range tests {
		actual := test.input.String()
		if actual != test.expected {
			t.Errorf("test %v: expected '%v', but got '%v'", i, test.expected, actual)
		}
	}
}
