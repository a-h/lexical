package input

import (
	"io"
	"testing"
)

func TestStringRuneReader(t *testing.T) {
	sr := &StringRuneReader{
		position: 0,
		s:        "ABC爱得林",
	}

	a, _, err := sr.ReadRune()
	expect(a, 'A', t)
	b, _, err := sr.ReadRune()
	expect(b, 'B', t)
	c, _, err := sr.ReadRune()
	expect(c, 'C', t)
	ai, _, err := sr.ReadRune()
	expect(ai, '爱', t)
	de, _, err := sr.ReadRune()
	expect(de, '得', t)
	lin, _, err := sr.ReadRune()
	expect(lin, '林', t)
	_, _, err = sr.ReadRune()
	if err != io.EOF {
		t.Error("expected EOF, but didn't get it.")
	}
}

func expect(actual, expected rune, t *testing.T) {
	if actual != expected {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestStringRuneReaderReadPastEnd(t *testing.T) {
	sr := &StringRuneReader{
		position: 0,
		s:        "A",
	}

	a, _, err := sr.ReadRune()
	expect(a, 'A', t)
	if err != nil {
		t.Error("unexpected error")
	}
	// We've hit the end.
	_, _, err = sr.ReadRune()
	if err != io.EOF {
		t.Error("expected EOF, but didn't get it.")
	}
	// Even though we got an error, ignore it and try again.
	_, _, err = sr.ReadRune()
	if err != io.EOF {
		t.Error("expected EOF, but didn't get it.")
	}
}

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
		s := NewFromString(test.input)

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
	s := NewFromString("ABCDEFG")

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
	s := NewFromString("你叫什么name？")

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
	s := NewFromString("ABCDEFG")

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

func TestStreamPeek(t *testing.T) {
	s := NewFromString("ABCDEFG")

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

func TestStreamPeekN(t *testing.T) {
	s := NewFromString("ABCDEFG")

	startPosition := s.Current
	startRune := s.CurrentRune

	peekedString, _ := s.PeekN(3)
	if peekedString != "ABC" {
		t.Errorf("Expected to peek 'ABC', but got '%v'", peekedString)
	}

	if s.Current != startPosition {
		t.Error("Peeking shouldn't modify the current position")
	}

	if s.CurrentRune != startRune {
		t.Error("Peeking shouldn't change the current rune")
	}
}

func TestStreamCollect(t *testing.T) {
	s := NewFromString("ABCDEFG")

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
	s := NewFromString("ABCDEFG")

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

func BenchmarkStreamAdvance(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		s := NewFromString("ABCDEFG")
		s.Advance()
	}
}

func BenchmarkStreamRetreat(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		s := NewFromString("ABCDEFG")
		s.Advance()
		s.Retreat()
	}
}

func BenchmarkStreamPeek(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		s := NewFromString("ABCDEFG")
		s.Peek()
	}
}

func BenchmarkStreamPosition(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		s := NewFromString("ABCDEFG")
		s.Position()
	}
}
