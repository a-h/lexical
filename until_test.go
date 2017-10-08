package lexical

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"unicode"
)

func TestUntil(t *testing.T) {
	bs := bytes.NewBufferString("Abc Def")
	sr := bufio.NewReader(bs)
	var wordReader func(lex *Lexer) StateFunction
	wordReader = func(lex *Lexer) StateFunction {
		// Read until a space and emit the word.
		lex.AdvanceUntilRune(' ')
		item := lex.Emit(ItemTypeWord)
		if item.Value != "Abc" {
			t.Errorf("expected to emit 'Abc', but actually emitted '%v': %v", item.Value, lex)
		}
		lex.AdvanceUntil(func(r rune) bool { return !unicode.IsSpace(r) })
		item = lex.Emit(ItemTypeSpace)
		if item.Value != " " {
			t.Errorf("expected to emit a space, but actually emitted '%v': %v", item.Value, lex)
		}
		return nil
	}
	l := NewLexer("words", sr, wordReader)

	actual := []Item{}
	for item := range l.Items {
		actual = append(actual, item)
	}
}

func TestUntilAdvanced(t *testing.T) {
	expected := []Item{
		Item{
			Type:  ItemTypeWord,
			Value: "Abc",
		},
		Item{
			Type:  ItemTypeSpace,
			Value: " ",
		},
		Item{
			Type:  ItemTypeWord,
			Value: "Def",
		},
		Item{
			Type:  ItemTypeEOF,
			Value: "",
		},
	}

	bs := bytes.NewBufferString("Abc Def")
	sr := bufio.NewReader(bs)
	var wordReader func(lex *Lexer) StateFunction
	wordReader = func(lex *Lexer) StateFunction {
		// Read until a space and emit the word.
		err := lex.AdvanceUntilRune(' ')
		lex.Emit(ItemTypeWord)
		if err == io.EOF {
			lex.Emit(ItemTypeEOF)
			return nil
		}
		// Ignore everything that's not a space.
		err = lex.AdvanceUntil(func(r rune) bool { return !unicode.IsSpace(r) })
		if err == nil {
			lex.Emit(ItemTypeSpace)
		}
		if err == io.EOF {
			lex.Emit(ItemTypeEOF)
			return nil
		}
		return wordReader
	}
	l := NewLexer("words", sr, wordReader)

	actual := []Item{}
	for item := range l.Items {
		actual = append(actual, item)
	}

	if len(actual) != len(expected) {
		t.Errorf("expected %v results, got %v", len(expected), len(actual))
		return
	}

	for i, a := range actual {
		e := expected[i]
		if !e.Eq(a) {
			t.Errorf("for index %v, expected '%v', but got '%v'", i, e, a)
			if e.Value != a.Value {
				t.Errorf("for index %v, expected value '%v', but got '%v'", i, e.Value, a.Value)
			}
			if e.Type != a.Type {
				t.Errorf("for index %v, expected item type '%v', but got '%v'", i, e.Type, a.Type)
			}
		}
	}
}
