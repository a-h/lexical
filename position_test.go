package lexical

import (
	"bufio"
	"bytes"
	"testing"
)

func TestLexerPositionAdvance(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Position
	}{
		{
			name:  "single line",
			input: "12345",
			expected: []Position{
				NewPosition(1, 1),
				NewPosition(1, 2),
				NewPosition(1, 3),
				NewPosition(1, 4),
				NewPosition(1, 5),
			},
		},
		{
			name:  "two lines",
			input: "a\nb\n",
			expected: []Position{
				NewPosition(1, 1), // a
				NewPosition(2, 0), // \n
				NewPosition(2, 1), // b
				NewPosition(3, 0), // \n
			},
		},
		{
			name:  "three lines",
			input: "11\n212\n31",
			expected: []Position{
				NewPosition(1, 1),
				NewPosition(1, 2),
				NewPosition(2, 0),
				NewPosition(2, 1),
				NewPosition(2, 2),
				NewPosition(2, 3),
				NewPosition(3, 0),
				NewPosition(3, 1),
				NewPosition(3, 2),
			},
		},
		{
			name:  "four lines",
			input: "\n\n\n\n",
			expected: []Position{
				NewPosition(2, 0),
				NewPosition(3, 0),
				NewPosition(4, 0),
				NewPosition(5, 0),
			},
		},
		{
			name:  "windows line break",
			input: "a\r\nb",
			expected: []Position{
				NewPosition(1, 1), // a
				NewPosition(1, 1), // \r
				NewPosition(2, 0), // \n
				NewPosition(2, 1), // b
			},
		},
	}

	advanceOperation := func(lex *Lexer) (rune, error) { return lex.Advance() }

	for _, test := range tests {
		advanceCount := 0
		actual := testPosition(test.name, test.input, advanceCount, advanceOperation, t)

		if len(test.expected) != len(actual) {
			t.Errorf("name: '%s': expected %d positions, but got %d positions", test.name, len(test.expected), len(actual))
		}

		// Check the positions.
		for i, e := range test.expected {
			a := actual[i]
			if !e.Eq(a) {
				t.Errorf("name '%s': index %d, expected position %v, but got %v", test.name, i, e.String(), a.String())
			}
		}
	}
}

// testPosition tests applies the operation to the lexer and checks the results. The advanceCount does some inital setup on the lexer.
func testPosition(name string, input string, advanceCount int, operation func(lex *Lexer) (rune, error), t *testing.T) []Position {
	bs := bytes.NewBufferString(input)
	sr := bufio.NewReader(bs)

	actual := make([]Position, 0)
	initialSetup := true

	sf := func(lex *Lexer) StateFunction {
		if initialSetup {
			for i := 0; i < advanceCount; i++ {
				lex.Advance()
			}
			initialSetup = false
		}

		for {
			if _, err := operation(lex); err != nil {
				break
			}
			actual = append(actual, lex.Position)
		}
		return nil
	}
	l := NewLexer(name, sr, sf)

	// Try and read the items, but there aren't any.
	for _ = range l.Items {
	}

	return actual
}

func TestLexerPositionRetreat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Position
	}{
		{
			name:  "single line",
			input: "12345",
			expected: []Position{
				NewPosition(1, 5),
				NewPosition(1, 4),
				NewPosition(1, 3),
				NewPosition(1, 2),
				NewPosition(1, 1),
			},
		},
		{
			name:  "two lines",
			input: "a\nb",
			expected: []Position{
				// b isn't included
				NewPosition(2, 1), // \n
				NewPosition(1, 1), //
				NewPosition(1, 0), // start
			},
		},
		{
			name:  "windows line break",
			input: "a\r\nb",
			expected: []Position{
				// b isn't included
				NewPosition(2, 1), // \n, we still stay on the line
				NewPosition(1, 1), // \r
				NewPosition(1, 1), //
				NewPosition(1, 0), // start
			},
		},
		{
			name:  "groups of 3",
			input: "123\n456\n789",
			expected: []Position{
				// 9 isn't included
				NewPosition(3, 3), // 8
				NewPosition(3, 2), // 7

				NewPosition(3, 1), // \n

				NewPosition(2, 3), // 6
				NewPosition(2, 2), // 5
				NewPosition(2, 1), // 4

				NewPosition(2, 0), // \n

				NewPosition(1, 3), // 3
				NewPosition(1, 2), // 2
				NewPosition(1, 1), // 1

				NewPosition(1, 0), // start
			},
		},
	}

	retreatOperation := func(lex *Lexer) (rune, error) { return lex.Retreat() }

	for _, test := range tests {
		actual := testPosition(test.name, test.input, len(test.input)+1, retreatOperation, t)

		if len(test.expected) != len(actual) {
			t.Errorf("name: '%s': expected %d positions, but got %d positions", test.name, len(test.expected), len(actual))
			break
		}

		// Check the positions.
		for i, e := range test.expected {
			a := actual[i]
			if !e.Eq(a) {
				t.Errorf("name '%s': index %d, expected position %v, but got %v", test.name, i, e.String(), a.String())
			}
		}
	}
}

func TestPositionString(t *testing.T) {
	tests := []struct {
		input    Position
		expected string
	}{
		{
			input:    NewPosition(1, 1),
			expected: "Line: 1, Col: 1",
		},
		{
			input:    NewPosition(2, 1),
			expected: "Line: 2, Col: 1",
		},
		{
			input:    NewPosition(2, 3),
			expected: "Line: 2, Col: 3",
		},
	}

	for i, test := range tests {
		actual := test.input.String()
		if actual != test.expected {
			t.Errorf("index %d, for line %v and col %v expected string '%v', but got '%v'", i,
				test.input.Line, test.input.Col, test.expected, actual)
		}
	}
}

func TestPositionEquals(t *testing.T) {
	tests := []struct {
		a        Position
		b        Position
		expected bool
	}{
		{
			a:        NewPosition(1, 1),
			b:        NewPosition(1, 1),
			expected: true,
		},
		{
			a:        NewPosition(2, 1),
			b:        NewPosition(2, 1),
			expected: true,
		},
		{
			a:        NewPosition(1, 2),
			b:        NewPosition(1, 2),
			expected: true,
		},
		{
			a:        NewPosition(3, 2),
			b:        NewPosition(3, 1),
			expected: false,
		},
		{
			a:        NewPosition(3, 1),
			b:        NewPosition(4, 1),
			expected: false,
		},
	}

	for i, test := range tests {
		actual := test.a.Eq(test.b)
		if actual != test.expected {
			t.Errorf("index %d, expected position %v eq %v to be %v, but was %v", i, test.a, test.b, actual, test.expected)
		}
	}
}

func TestPositionAdvance(t *testing.T) {
	p := NewPosition(1, 1)
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 2 {
		t.Errorf("expected a newline to advance the position from one to two")
	}
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 3 {
		t.Errorf("expected a newline to advance the position from two to three")
	}
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 4 {
		t.Errorf("expected a newline to advance the position from three to four")
	}
}

func TestPositionRetreat(t *testing.T) {
	p := NewPosition(1, 1)
	p.Advance('\n')
	comparePosition(NewPosition(2, 0), p, t)
	p.Advance('a')
	comparePosition(NewPosition(2, 1), p, t)
	p.Advance('b')
	comparePosition(NewPosition(2, 2), p, t)
	p.Advance('\n')
	comparePosition(NewPosition(3, 0), p, t)
	p.Advance('b')
	comparePosition(NewPosition(3, 1), p, t)
	p.Retreat('\n')
	comparePosition(NewPosition(2, 2), p, t)
	p.Retreat('b')
	comparePosition(NewPosition(2, 1), p, t)
	p.Retreat('\n')
	comparePosition(NewPosition(1, 1), p, t)
	p.Retreat('a')
	comparePosition(NewPosition(1, 0), p, t)
}

func comparePosition(expected, actual Position, t *testing.T) {
	if !expected.Eq(actual) {
		t.Errorf("expected %v, but got %v", expected.String(), actual.String())
	}
}
