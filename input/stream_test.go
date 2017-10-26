package input

import (
	"testing"
)

func TestStreamAdvance(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "Abc Bcd",
		},
		{
			input: "你好",
		},
	}

	for _, test := range tests {
		s := NewFromString("Advance Test", test.input)

		for i, er := range test.input {
			ar, err := s.Advance()
			if err != nil {
				t.Errorf("for input '%v', had fatal error at index %v: %v", test.input, i, err)
			}
			if ar != er {
				t.Errorf("for input '%v', failed to advance to rune %v. Expected %v, got %v.", test.input, i, string(er), string(ar))
			}
		}
	}
}

func expectRune(s *Stream, test func() (rune, error), expected rune, t *testing.T, name string) {
	actual, err := test()
	if err != nil {
		t.Errorf("%v: expected rune: '%v', but got error: %v", name, string(expected), err)
	}
	if actual != expected {
		t.Errorf("%v: expected rune: '%v', but got '%v'", name, string(expected), string(actual))
	}
}

func TestStreamRetreat(t *testing.T) {
	s := NewFromString("Retreat Test", "ABCDEFG")

	expectRune(s, s.Advance, 'A', t, "1")
	expectRune(s, s.Advance, 'B', t, "2")
	expectRune(s, s.Advance, 'C', t, "3")
	expectRune(s, s.Advance, 'D', t, "4")
	expectRune(s, s.Retreat, 'C', t, "5")
	expectRune(s, s.Retreat, 'B', t, "6")
	expectRune(s, s.Retreat, 'A', t, "7")
	expectRune(s, s.Advance, 'B', t, "8")
	expectRune(s, s.Advance, 'C', t, "9")
	expectRune(s, s.Advance, 'D', t, "10")
	expectRune(s, s.Advance, 'E', t, "11")
	expectRune(s, s.Advance, 'F', t, "12")
	expectRune(s, s.Advance, 'G', t, "13")
}

func TestStreamAdvanceRetreat_UTF8(t *testing.T) {
	s := NewFromString("words", "你叫什么name？")

	expectRune(s, s.Advance, '你', t, "1")
	expectRune(s, s.Advance, '叫', t, "2")
	expectRune(s, s.Advance, '什', t, "3")
	expectRune(s, s.Advance, '么', t, "4")
	expectRune(s, s.Advance, 'n', t, "5")
	expectRune(s, s.Advance, 'a', t, "6")
	expectRune(s, s.Advance, 'm', t, "7")
	expectRune(s, s.Advance, 'e', t, "8")
	expectRune(s, s.Advance, '？', t, "9")
	expectRune(s, s.Retreat, 'e', t, "10")
	expectRune(s, s.Retreat, 'm', t, "11")
	expectRune(s, s.Retreat, 'a', t, "12")
	expectRune(s, s.Retreat, 'n', t, "13")
	expectRune(s, s.Retreat, '么', t, "14")
	expectRune(s, s.Retreat, '什', t, "15")
	expectRune(s, s.Retreat, '叫', t, "16")
	expectRune(s, s.Retreat, '你', t, "17")
	expectRune(s, s.Advance, '叫', t, "18")
}

func TestStreamRetreat_CannotReadBeforeStartOfStream(t *testing.T) {
	s := NewFromString("words", "ABCDEFG")

	// Read the first two runes.
	expectRune(s, s.Advance, 'A', t, "advance to a")
	expectRune(s, s.Advance, 'B', t, "advance to b")

	// Retreat once.
	expectRune(s, s.Retreat, 'A', t, "retreat back to a")

	// Expect an empty rune because we're right back at the start of the stream.
	_, err := s.Retreat()
	if err != ErrStartOfFile {
		t.Errorf("expected to be back at the start of the stream")
	}

	// But can't go past the start of the stream.
	r, err := s.Retreat()
	if err == nil {
		t.Errorf("it should not be possible to retreat past the start of a stream, but got rune '%v'", string(r))
	}
}

func TestStreamLeftAndRight(t *testing.T) {
	tests := []struct {
		input  string
		middle int
		left   string
		right  string
	}{
		{
			input:  "abcd",
			middle: 0,
			left:   "",
			right:  "abcd",
		},
		{
			input:  "abcd",
			middle: 1,
			left:   "a",
			right:  "bcd",
		},
		{
			input:  "abcd",
			middle: 2,
			left:   "ab",
			right:  "cd",
		},
		{
			input:  "abcd",
			middle: 3,
			left:   "abc",
			right:  "d",
		},
		{
			input:  "abcd",
			middle: 4,
			left:   "abcd",
			right:  "",
		},
		{
			input:  "abcd",
			middle: 5,
			left:   "abcd",
			right:  "",
		},
	}

	for _, test := range tests {
		l := getLeft([]rune(test.input), test.middle)
		if string(l) != test.left {
			t.Errorf("for input '%v', at middle %v, expected left of '%v' but got '%v'", test.input, test.middle, test.left, string(l))
		}
		r := getRight([]rune(test.input), test.middle)
		if string(r) != test.right {
			t.Errorf("for input '%v', at middle %v, expected right of '%v' but got '%v'", test.input, test.middle, test.right, string(l))
		}
	}
}

func TestStreamPeek(t *testing.T) {
	s := NewFromString("words", "ABCDEFG")

	startPosition := s.Current
	startRune := s.CurrentRune

	peekedRune, _ := s.Peek()
	if peekedRune != 'A' {
		t.Errorf("Expected to peek 'A', but got '%v'", peekedRune)
	}

	if s.Current != startPosition {
		t.Error("Peeking shouldn't modify the current position")
	}

	if s.CurrentRune != startRune {
		t.Error("Peeking shouldn't change the current rune")
	}
}

func TestStreamCollect(t *testing.T) {
	s := NewFromString("words", "ABCDEFG")

	s.Advance() // A
	s.Advance() // B
	s.Advance() // C

	abc := s.Collect()
	if s.Index() != 3 {
		t.Errorf("Expected position 3, but got %d", s.Index())
	}
	if abc != "ABC" {
		t.Errorf("Expected to collect 'ABC', but got '%v'. Stream was %v", abc, s)
	}

	s.Advance() // D
	s.Advance() // E
	s.Advance() // F
	s.Advance() // G

	defg := s.Collect()
	if s.Index() != 7 {
		t.Errorf("Expected position 7, but got %d", s.Index())
	}
	if defg != "DEFG" {
		t.Errorf("Expected to collect 'DEFG', but got '%v'", defg)
	}
}

func TestStreamAdvanceCollectRetreat(t *testing.T) {
	s := NewFromString("words", "ABCDEFG")

	expectRune(s, s.Advance, 'A', t, "1")
	expectRune(s, s.Advance, 'B', t, "2")
	expectRune(s, s.Advance, 'C', t, "3")

	expectRune(s, s.Retreat, 'B', t, "4")

	ab := s.Collect()
	if ab != "AB" {
		t.Errorf("Expected to collect 'AB', but got '%v'", ab)
	}

	expectRune(s, s.Advance, 'C', t, "5")
	expectRune(s, s.Advance, 'D', t, "6")
	s.Advance() // E
	s.Advance() // F

	cdef := s.Collect()
	if cdef != "CDEF" {
		t.Errorf("Expected to collect 'CDEF', but got '%v'", cdef)
		t.Error(s)
	}
}