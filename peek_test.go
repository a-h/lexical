package lexical

import (
	"bufio"
	"bytes"
	"testing"
)

func TestPeek(t *testing.T) {
	bs := bytes.NewBufferString("Ab")
	sr := bufio.NewReader(bs)

	var a, aa, b, bb rune

	runeReader := func(lex *Lexer) StateFunction {
		var err error
		a, err = lex.Peek()
		if err != nil {
			t.Errorf("failed to peek the first rune with error: '%v'", err)
		}
		aa, err = lex.Advance()
		if err != nil {
			t.Errorf("failed to advance the first rune with error: '%v'", err)
		}
		b, err = lex.Peek()
		if err != nil {
			t.Errorf("failed to peek the second rune with error: '%v'", err)
		}
		bb, err = lex.Advance()
		if err != nil {
			t.Errorf("failed to advance the second rune with error: '%v'", err)
		}
		return nil
	}
	l := NewLexer("words", sr, runeReader)

	// Try and read the items, but there aren't any.
	for _ = range l.Items {
	}

	if a != 'A' {
		t.Errorf("failed to peek the first rune, expected 'A', but got '%v'", a)
	}
	if aa != 'A' {
		t.Errorf("failed to read the first rune, expected 'A', but got '%v'", a)
	}
	if b != 'b' {
		t.Errorf("failed to peek the first rune, expected 'b', but got '%v'", b)
	}
	if bb != 'b' {
		t.Errorf("failed to read the first rune, expected 'b', but got '%v'", b)
	}
}
