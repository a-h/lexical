package input

import (
	"testing"
)

func TestPositionAdvance(t *testing.T) {
	p := NewPosition(1, 0)
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 2 {
		t.Errorf("expected a newline to advance the position from one to two")
	}
	if p.Col != 1 {
		t.Errorf("expected a newline to bounce to the zeroth character on the next line")
	}
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 3 {
		t.Errorf("expected a newline to advance the position from two to three")
	}
	if p.Col != 1 {
		t.Errorf("expected a newline to bounce to the zeroth character on the next line")
	}
	p.Advance('\n')
	p.Advance('a')
	if p.Line != 4 {
		t.Errorf("expected a newline to advance the position from three to four")
	}
	if p.Col != 1 {
		t.Errorf("expected a newline to bounce to the zeroth character on the next line")
	}
}

func TestPositionRetreatOverNewLine(t *testing.T) {
	p := NewPosition(1, 0)
	p.Advance('a') // Index 0
	comparePosition(NewPosition(1, 1), p, t)
	p.Advance('\n') // Index 1
	comparePosition(NewPosition(2, 0), p, t)
	p.Retreat('\n') // Index 0
	comparePosition(NewPosition(1, 1), p, t)
	p.Retreat('a') // Start of file
	comparePosition(NewPosition(1, 0), p, t)
	if p.Index != -1 {
		t.Errorf("should be at start of file")
	}
}

func TestPositionRetreat(t *testing.T) {
	p := NewPosition(1, 0) // Start of file.
	p.Advance('a')
	comparePosition(NewPosition(1, 1), p, t)
	p.Advance('b')
	comparePosition(NewPosition(1, 2), p, t)
	p.Advance('\n')
	comparePosition(NewPosition(2, 0), p, t)
	p.Advance('c')
	comparePosition(NewPosition(2, 1), p, t)
	p.Retreat('c')
	comparePosition(NewPosition(2, 0), p, t)
	p.Retreat('\n')
	comparePosition(NewPosition(1, 2), p, t)
	p.Retreat('b')
	comparePosition(NewPosition(1, 1), p, t)
	p.Retreat('a')
	comparePosition(NewPosition(1, 0), p, t)
}

func comparePosition(expected, actual Position, t *testing.T) {
	if !expected.Eq(actual) {
		t.Errorf("expected %v, but got %v", expected.String(), actual.String())
	}
}
func TestPositionAdvanceRetreatNewLine(t *testing.T) {
	actual := NewPosition(1, 0)
	actual.Advance('\n')
	expected := NewPosition(2, 0)
	if !expected.Eq(actual) {
		t.Errorf("advance 0: '\\n': expected %v, but got %v", expected.String(), actual.String())
	}
	actual.Advance('\n')
	expected = NewPosition(3, 0)
	if !expected.Eq(actual) {
		t.Errorf("advance 1: '\\n': expected %v, but got %v", expected.String(), actual.String())
	}
	actual.Retreat('\n')
	expected = NewPosition(2, 0)
	if !expected.Eq(actual) {
		t.Errorf("retreat 2: '\\n': expected %v, but got %v", expected.String(), actual.String())
	}
	actual.Retreat('\n')
	expected = NewPosition(1, 0)
	if !expected.Eq(actual) {
		t.Errorf("retreat 3: '\\n': expected %v, but got %v", expected.String(), actual.String())
	}
}

func TestPositionAdvanceRetreat(t *testing.T) {
	input := "\nab\nc\nd"
	expectedAdvances := []Position{
		NewPosition(2, 0), // \n
		NewPosition(2, 1), // a
		NewPosition(2, 2), // b
		NewPosition(3, 0), // \n
		NewPosition(3, 1), // c
		NewPosition(4, 0), // \n
		NewPosition(4, 1), // d
	}
	// Advance.
	actual := NewPosition(1, 0)
	for i, r := range input {
		actual.Advance(r)
		expected := expectedAdvances[i]
		if !expected.Eq(actual) {
			t.Errorf("advance %d (%s): expected %v, but got %v", i, string(r), expected.String(), actual.String())
		}
	}
	if !actual.Eq(NewPosition(4, 1)) {
		t.Errorf("after reading everything, expected %v, got %v", NewPosition(4, 1), actual)
	}

	// Unread everything.
	input = "d\nc\nba\n"
	expectedRetreats := []Position{
		NewPosition(4, 0), // \d
		NewPosition(3, 1), // \n
		NewPosition(3, 0), // \c
		NewPosition(2, 2), // \n
		NewPosition(2, 1), // b
		NewPosition(2, 0), // a
		NewPosition(1, 0), // \n
	}
	// Retreat.
	for i := 0; i > 0; i-- {
		r := rune(input[i])
		actual.Retreat(r)
		expected := expectedRetreats[i]
		if !expected.Eq(actual) {
			t.Errorf("retreat %d (%s): expected %v, but got %v", i, string(r), expected.String(), actual.String())
		}
	}
}

func TestPositionString(t *testing.T) {
	p := NewPosition(1, 1)
	if p.String() != "Line: 1, Col: 1" {
		t.Errorf("Expected 'Line: 1, Col: 1', but got %v", p.String())
	}
	p.Advance('a')
	if p.String() != "Line: 1, Col: 2" {
		t.Errorf("Expected 'Line: 1, Col: 2', but got %v", p.String())
	}
	p.Advance('\n')
	if p.String() != "Line: 2, Col: 0" {
		t.Errorf("Expected 'Line: 2, Col: 0', but got %v", p.String())
	}
	p.Advance('b')
	if p.String() != "Line: 2, Col: 1" {
		t.Errorf("Expected 'Line: 2, Col: 1', but got %v", p.String())
	}
}

func BenchmarkPosition(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		p := NewPosition(1, 1)
		p.Advance('a')
		p.Advance('\n')
		p.Retreat('\n')
	}
}
