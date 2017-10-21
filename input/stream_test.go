package input

import (
	"bufio"
	"bytes"
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
		bs := bytes.NewBufferString(test.input)
		sr := bufio.NewReader(bs)
		s := New("Advance Test", sr)

		for i, er := range test.input {
			ar, err := s.Advance()
			if err != nil {
				t.Errorf("for input '%v', had fatal error at index %v: %v", test.input, i, err)
			}
			if ar != er {
				t.Errorf("for input '%v', failed to advance to rune %v. Expected %d, got %d.", test.input, i, er, ar)
			}
		}
	}
}

func expectRune(test func() (rune, error), expected rune, t *testing.T, name string) {
	actual, err := test()
	if err != nil {
		t.Errorf("%v: expected rune: '%v', but got error: %v", name, string(expected), err)
	}
	if actual != expected {
		t.Errorf("%v: expected rune: '%v', but got '%v'", name, string(expected), string(actual))
	}
}

func TestStreamRetreat(t *testing.T) {
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	s := New("Retreat Test", sr)

	expectRune(s.Advance, 'A', t, "1")
	expectRune(s.Advance, 'B', t, "2")
	expectRune(s.Advance, 'C', t, "3")
	expectRune(s.Advance, 'D', t, "4")
	expectRune(s.Retreat, 'C', t, "5")
	expectRune(s.Retreat, 'B', t, "6")
	expectRune(s.Retreat, 'A', t, "7")
	expectRune(s.Advance, 'B', t, "8")
	expectRune(s.Advance, 'C', t, "9")
	expectRune(s.Advance, 'D', t, "10")
	expectRune(s.Advance, 'E', t, "11")
	expectRune(s.Advance, 'F', t, "12")
	expectRune(s.Advance, 'G', t, "13")
}

func TestStreamAdvanceRetreat_UTF8(t *testing.T) {
	bs := bytes.NewBufferString("你叫什么name？")
	sr := bufio.NewReader(bs)

	s := New("words", sr)
	expectRune(s.Advance, '你', t, "1")
	expectRune(s.Advance, '叫', t, "2")
	expectRune(s.Advance, '什', t, "3")
	expectRune(s.Advance, '么', t, "4")
	expectRune(s.Advance, 'n', t, "5")
	expectRune(s.Advance, 'a', t, "6")
	expectRune(s.Advance, 'm', t, "7")
	expectRune(s.Advance, 'e', t, "8")
	expectRune(s.Advance, '？', t, "9")
	expectRune(s.Retreat, 'e', t, "10")
	expectRune(s.Retreat, 'm', t, "11")
	expectRune(s.Retreat, 'a', t, "12")
	expectRune(s.Retreat, 'n', t, "13")
	expectRune(s.Retreat, '么', t, "14")
	expectRune(s.Retreat, '什', t, "15")
	expectRune(s.Retreat, '叫', t, "16")
	expectRune(s.Retreat, '你', t, "17")
	expectRune(s.Advance, '叫', t, "18")
}

func TestStreamRetreat_CannotReadBeforeStartOfStream(t *testing.T) {
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	s := New("words", sr)

	// Read the first two runes.
	expectRune(s.Advance, 'A', t, "advance to a")
	expectRune(s.Advance, 'B', t, "advance to b")

	// Retreat once.
	expectRune(s.Retreat, 'A', t, "retreat back to a")

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
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	s := New("words", sr)

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
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	s := New("words", sr)

	s.Advance() // A
	s.Advance() // B
	s.Advance() // C

	abc := s.Collect()
	if abc != "ABC" {
		t.Errorf("Expected to collect 'ABC', but got '%v'", abc)
	}

	s.Advance() // D
	s.Advance() // E
	s.Advance() // F
	s.Advance() // G

	defg := s.Collect()
	if defg != "DEFG" {
		t.Errorf("Expected to collect 'DEFG', but got '%v'", defg)
	}
}
