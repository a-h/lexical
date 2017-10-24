package input

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
)

// Stream defines a lexical scanner over a stream.
type Stream struct {
	// Name is used for error reports.
	Name string
	// Input holds the Reader being scanned.
	Input *bufio.Reader
	// Buffer is the space currently being searched for tokens to avoid seeking the input stream.
	// When a token match is found, the buffer is emptied.
	Buffer []rune
	// readUntil is the furthest point in the input we've read to. If this is ahead of Current, then we
	// have data in the Buffer which can be read.
	readUntil int64
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
	return fmt.Sprintf("%v: Current Rune: '%v', Start of Token Position: %v, Current Position: %v, Forward Buffer Size: %v, Current Buffer: '%v'", l.Name, string(l.CurrentRune), l.Start, l.Current, len(l.Buffer), string(l.Buffer))
}

// New creates a new parser input from a buffered reader.
func New(name string, input *bufio.Reader) *Stream {
	return &Stream{
		Name:      name,
		Input:     input,
		Buffer:    make([]rune, 0),
		Current:   -1,
		readUntil: -1,
		position:  NewPosition(1, 0),
	}
}

// NewFromString creates a new parser input from an input string.
func NewFromString(name string, input string) *Stream {
	bs := bytes.NewBufferString(input)
	sr := bufio.NewReader(bs)
	return New(name, sr)
}

// Collect returns the value of the consumed buffer and updates the position of the stream to the current
// position.
func (l *Stream) Collect() string {
	// Emit the token and update the position of the lexer against the input stream.
	// Returning the item helps with unit testing.
	length := int(l.Current - l.Start + 1)
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

// Advance reads a rune from the Input and sets the current position.
func (l *Stream) Advance() (rune, error) {
	// Check to see whether we already have it in the buffer, if so, read it from there.
	if l.Current+1 <= l.readUntil {
		l.Current++
		l.CurrentRune = l.Buffer[l.Current-l.Start]
		l.position.Advance(l.CurrentRune)
		return l.CurrentRune, nil
	}

	r, _, err := l.Input.ReadRune()
	l.Buffer = append(l.Buffer, r)
	l.Current++
	l.readUntil = l.Current
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

// ErrStartOfFile is the error used when we've retreated to the start of the file and can't
// retreat any further.
var ErrStartOfFile = errors.New("SOF")

// Retreat steps back a rune.
func (l *Stream) Retreat() (rune, error) {
	newPos := l.Current - 1
	if newPos <= -1 {
		l.Current = -1
		l.CurrentRune = 0x0
		l.position.Retreat(l.CurrentRune)
		return 0x0, ErrStartOfFile
	}
	l.Current = newPos
	l.CurrentRune = l.Buffer[l.Current-l.Start]
	l.position.Retreat(l.CurrentRune)
	return l.CurrentRune, nil
}

// Position returns the current position within the stream.
func (l *Stream) Position() (line, column int) {
	return l.position.Line, l.position.Col
}

// Index returns the current index position within the stream.
func (l *Stream) Index() int64 {
	return l.Current
}
