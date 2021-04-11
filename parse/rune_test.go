package parse

import (
	"testing"

	"github.com/a-h/lexical/input"
)

func TestRune(t *testing.T) {
	tests := []struct {
		input         string
		r             rune
		expected      bool
		expectedIndex int64
	}{
		{
			input:         "A",
			r:             'A',
			expected:      true,
			expectedIndex: 1,
		},
		{
			input:         "b",
			r:             'A',
			expected:      false,
			expectedIndex: 0,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		parser := Rune(test.r)
		result := parser(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if test.expected && result.Item != test.r {
			t.Errorf("test %v: for input '%v' expected to capture '%v' but got '%v'", i, test.input, test.r, result.Item)
		}
		if test.expectedIndex != pi.Index() {
			t.Errorf("test %v: for input '%v' expected index %d, got %d", i, test.input, test.expectedIndex, pi.Index())
		}
	}
}

func TestRunePosition(t *testing.T) {
	pi := input.NewFromString("ABCABC")
	pi.Advance() // A
	pi.Advance() // B
	pi.Advance() // C
	result := Rune('B')(pi)
	if result.Success {
		t.Errorf("did not expect 'B'")
	}
	if pi.Index() != 3 {
		t.Errorf("failed rune parse should rollback position")
	}
	result = Rune('A')(pi)
	if !result.Success {
		t.Errorf("expected 'A'")
	}
	if pi.Index() != 4 {
		t.Errorf("succesful parsing advanced position")
	}
}

func BenchmarkRune(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		parser := Rune('A')
		parser(input.NewFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}
}
