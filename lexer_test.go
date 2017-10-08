package lexical

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"unicode"
)

func TestThatTheLexerCanAdvanceThroughAReader(t *testing.T) {
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
		l := NewLexer("words", sr, End)

		for i, er := range test.input {
			ar, err := l.Advance()
			if err != nil {
				t.Errorf("for input '%v', had fatal error at index %v: %v", test.input, i, err)
			}
			if ar != er {
				t.Errorf("for input '%v', failed to advance to rune %v. Expected %d, got %d.", test.input, i, er, ar)
			}
		}
	}
}

func TestThatTheLexerCanRetreatThroughAReader(t *testing.T) {
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)

	l := NewLexer("words", sr, End)
	expectRune(l.Advance, 'A', t, "1")
	expectRune(l.Advance, 'B', t, "2")
	expectRune(l.Advance, 'C', t, "3")
	expectRune(l.Advance, 'D', t, "4")
	expectRune(l.Retreat, 'C', t, "5")
	expectRune(l.Retreat, 'B', t, "6")
	expectRune(l.Retreat, 'A', t, "7")
	expectRune(l.Advance, 'B', t, "8")
	expectRune(l.Advance, 'C', t, "9")
	expectRune(l.Advance, 'D', t, "10")
	expectRune(l.Advance, 'E', t, "11")
	expectRune(l.Advance, 'F', t, "12")
	expectRune(l.Advance, 'G', t, "13")
}

func TestThatTheLexerCanRetreatThroughAMixedReader(t *testing.T) {
	bs := bytes.NewBufferString("你叫什么name？")
	sr := bufio.NewReader(bs)

	l := NewLexer("words", sr, End)
	expectRune(l.Advance, '你', t, "1")
	expectRune(l.Advance, '叫', t, "2")
	expectRune(l.Advance, '什', t, "3")
	expectRune(l.Advance, '么', t, "4")
	expectRune(l.Advance, 'n', t, "5")
	expectRune(l.Advance, 'a', t, "6")
	expectRune(l.Advance, 'm', t, "7")
	expectRune(l.Advance, 'e', t, "8")
	expectRune(l.Advance, '？', t, "9")
	expectRune(l.Retreat, 'e', t, "10")
	expectRune(l.Retreat, 'm', t, "11")
	expectRune(l.Retreat, 'a', t, "12")
	expectRune(l.Retreat, 'n', t, "13")
	expectRune(l.Retreat, '么', t, "14")
	expectRune(l.Retreat, '什', t, "15")
	expectRune(l.Retreat, '叫', t, "16")
	expectRune(l.Retreat, '你', t, "17")
	expectRune(l.Advance, '叫', t, "18")
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

func TestThatTheLexerCannotRetreatBeforeTheStartOfAStream(t *testing.T) {
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	l := NewLexer("words", sr, End)

	// Read the first two runes.
	expectRune(l.Advance, 'A', t, "advance to a")
	expectRune(l.Advance, 'B', t, "advance to b")

	// Retreat once.
	expectRune(l.Retreat, 'A', t, "retreat back to a")

	// Expect an empty rune because we're right back at the start of the stream.
	_, err := l.Retreat()
	if err != ErrStartOfFile {
		t.Errorf("expected to be back at the start of the stream")
	}

	// But can't go past the start of the stream.
	r, err := l.Retreat()
	if err == nil {
		t.Errorf("it should not be possible to retreat past the start of a stream, but got rune '%v'", string(r))
	}
}

func End(l *Lexer) StateFunction {
	return nil
}

func TestSimpleCharacterLexer(t *testing.T) {
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
	l := NewLexer("words", sr, wordStateStart)

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
		if e != a {
			t.Errorf("for index %v, expected '%v', but got '%v'", i, e, a)
		}
	}
}

func wordStateStart(l *Lexer) StateFunction {
	for {
		r, err := l.Advance()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Emit(ItemTypeError)
			break
		}
		// Check the first character.

		// If it's a letter, we're in a word.
		if unicode.IsLetter(r) {
			// Since we're the start state, we should emit anything we've collected so far.
			if l.Current > l.Start {
				l.Emit(ItemTypeSpace)
			}

			return wordState
		}
	}
	if l.Current > l.Start {
		l.Emit(ItemTypeSpace)
	}
	l.Emit(ItemTypeEOF)
	return nil
}

const ItemTypeSpace ItemType = 10
const ItemTypeWord ItemType = 11

func wordState(l *Lexer) StateFunction {
	for {
		r, err := l.Advance()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Emit(ItemTypeError)
			break
		}

		// Everything until we hit a space is a word.
		if unicode.IsSpace(r) {
			l.Emit(ItemTypeWord)
			return wordStateStart
		}
	}
	if l.Current > l.Start {
		l.Emit(ItemTypeWord)
	}
	l.Emit(ItemTypeEOF)
	return nil
}

func TestThatTheLexerCannotGetStuckInALoop(t *testing.T) {
	bs := bytes.NewBufferString("abcdefg")
	sr := bufio.NewReader(bs)
	var wordReader func(lex *Lexer) StateFunction
	wordReader = func(lex *Lexer) StateFunction {
		return wordReader
	}
	l := NewLexer("empty", sr, wordReader)

	actual := []Item{}
	for item := range l.Items {
		actual = append(actual, item)
	}

	if len(actual) != 1 {
		t.Errorf("expected only an error item, but got %v items", len(actual))
		return
	}
	if actual[0].Type != ItemTypeError {
		t.Errorf("expected an error, but got an item of type %v", actual[0].Type)
	}
	if actual[0].Value != "lexer: stuck in a loop at position -1" {
		t.Errorf("unexpected error message '%v'", actual[0].Value)
	}
}

func TestLeftAndRight(t *testing.T) {
	tests := []struct {
		input  string
		middle int
		left   string
		right  string
	}{
		{
			input:  "abcd",
			middle: 2,
			left:   "ab",
			right:  "cd",
		},
		{
			input:  "abcd",
			middle: 4,
			left:   "abcd",
			right:  "",
		},
		{
			input:  "abcd",
			middle: 0,
			left:   "",
			right:  "abcd",
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
