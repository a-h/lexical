package input

import (
	"testing"
)

func TestStreamPositionAdvance(t *testing.T) {
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

	advanceOperation := func(s *Stream) (rune, error) { return s.Advance() }

	for _, test := range tests {
		actual := testPosition(test.input, 0, advanceOperation, t)

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

// testPosition tests applies the operation to the Stream and checks the results. The advanceCount does some inital setup on the Stream.
func testPosition(input string, advanceCount int, operation func(stream *Stream) (rune, error), t *testing.T) []Position {
	actual := make([]Position, 0)

	s := NewFromString(input)

	for i := 0; i < advanceCount; i++ {
		s.Advance()
	}

	for {
		if _, err := operation(s); err != nil {
			break
		}
		actual = append(actual, s.position)
	}

	return actual
}

func TestStreamPositionRetreat(t *testing.T) {
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

	retreatOperation := func(s *Stream) (rune, error) { return s.Retreat() }

	for _, test := range tests {
		actual := testPosition(test.input, len(test.input)+1, retreatOperation, t)

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

func TestStreamPositionFunction(t *testing.T) {
	s := NewFromString("abc\n123")

	line, col := s.Position()
	if line != 1 && col != 1 {
		t.Errorf("Expected line 1 and col 1, but got %v:%v", line, col)
	}
	s.Advance() // 'a'
	line, col = s.Position()
	if line != 1 && col != 2 {
		t.Errorf("Expected line 1 and col 2, but got %v:%v", line, col)
	}
	s.Advance() // 'b'
	line, col = s.Position()
	if line != 1 && col != 3 {
		t.Errorf("Expected line 1 and col 3, but got %v:%v", line, col)
	}
}
