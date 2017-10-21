package input

import (
	"testing"
)

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
