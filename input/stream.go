package input

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// Stream defines a lexical scanner over a stream.
type Stream struct {
	// Input holds the Reader being scanned.
	Input io.RuneReader
	// Buffer is the space currently being searched for tokens to avoid seeking the input stream.
	// When a token match is found, the buffer is emptied.
	Buffer []rune
	// Start represents the start position of the lexer against the Input, e.g.
	// if we've lexed 124 runes already, that's where we're starting from.
	Start int64
	// Current represents the current position of the lexer. We might have lexed
	// 124 runes and have read 3 more runes without yet emitting a token, so current would
	// be 124+3=127
	Current int64
	// CurrentRune is the current rune at the cursor.
	CurrentRune rune
	// Position is the current position within the file.
	position Position
}

func (l *Stream) String() string {
	return fmt.Sprintf("Current Rune: '%v', Start of Token Position: %v, Current Position: %v, Forward Buffer Size: %v, Current Buffer: '%v'", string(l.CurrentRune), l.Start, l.Current, len(l.Buffer), string(l.Buffer))
}

// New creates a new parser input from a buffered reader.
func New(input io.RuneReader) *Stream {
	return &Stream{
		Input:    input,
		Buffer:   make([]rune, 0),
		position: NewPosition(1, 0),
	}
}

// StringRuneReader allows a string to be read rune-by-rune. It allocates slightly less variables than
// NewBufferString or NewReader.
type StringRuneReader struct {
	position int
	s        string
}

// ReadRune reads a rune from the underlying string.
func (sr *StringRuneReader) ReadRune() (r rune, size int, err error) {
	r, size = utf8.DecodeRuneInString(sr.s[sr.position:])
	if size == 0 {
		err = io.EOF
	}
	sr.position += size
	return
}

// NewFromString creates a new parser input from an input string.
func NewFromString(input string) *Stream {
	return New(&StringRuneReader{
		position: 0,
		s:        input,
	})
}

// Collect returns the value of the consumed buffer and updates the position of the stream to the current
// position.
func (l *Stream) Collect() string {
	// Emit the token and update the position of the lexer against the input stream.
	// Returning the item helps with unit testing.
	length := int(l.Current - l.Start)
	left := getLeft(l.Buffer, length)
	right := getRight(l.Buffer, length)
	l.Start = l.Current
	l.Buffer = right
	return string(left)
}

func getLeft(runes []rune, length int) []rune {
	if length > len(runes) {
		return runes
	}
	return runes[:length]
}

func getRight(runes []rune, start int) []rune {
	if start >= len(runes) {
		return make([]rune, 0)
	}
	return runes[start:]
}

func fromBuffer(startOfBufferIndex int64, currentIndex int64, buffer []rune) (r rune, ok bool) {
	index := int(currentIndex-startOfBufferIndex) - 1
	if index < 0 {
		// Can't read before the buffer.
		return 0x0, false
	}
	if index >= len(buffer) {
		return 0x0, false
	}
	return buffer[index], true
}

// Advance reads a rune from the Input and sets the current position.
func (l *Stream) Advance() (r rune, err error) {
	// Check to see whether we already have it in the buffer, if so, read it from there.
	l.Current++

	r, ok := fromBuffer(l.Start, l.Current, l.Buffer)
	if !ok {
		r, _, err = l.Input.ReadRune()
		l.Buffer = append(l.Buffer, r)
	}

	l.CurrentRune = r
	l.position.Advance(l.CurrentRune)
	return r, err
}

// Peek reads a rune from the Input, then sets the current position back.
func (l *Stream) Peek() (rune, error) {
	r, err := l.Advance()
	if err != nil {
		return r, err
	}
	_, err = l.Retreat()
	if err != ErrStartOfFile {
		return r, err
	}
	return r, nil
}

// PeekN reads a number of runes from the Input, then sets the current position back.
func (l *Stream) PeekN(n int) (s string, err error) {
	rs := make([]rune, n)
	var advanced int
	for i := 0; i < n; i++ {
		rs[i], err = l.Advance()
		if err != nil {
			break
		}
		advanced++
	}

	s = string(rs)

	for i := 0; i < n; i++ {
		l.Retreat()
	}

	return
}

// ErrStartOfFile is the error used when we've retreated to the start of the file and can't
// retreat any further.
var ErrStartOfFile = errors.New("SOF")

// Retreat steps back a rune.
func (l *Stream) Retreat() (r rune, err error) {
	l.Current--

	r, ok := fromBuffer(l.Start, l.Current, l.Buffer)
	if !ok {
		l.CurrentRune = 0x0
		l.position.Retreat(l.CurrentRune)
		return 0x0, ErrStartOfFile
	}

	l.CurrentRune = r
	l.position.Retreat(l.CurrentRune)
	return r, err
}

// Position returns the current position within the stream.
func (l *Stream) Position() (line, column int) {
	return l.position.Line, l.position.Col
}

// Index returns the current index position within the stream.
func (l *Stream) Index() int64 {
	return l.Current
}
