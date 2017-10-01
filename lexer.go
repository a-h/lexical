package lexical

import (
	"bufio"
	"errors"
)

// Lexer defines a lexical scanner.
type Lexer struct {
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
	// Items is a channel where Items can be emitted.
	Items chan Item
}

// NewLexer creates a new Lexer.
func NewLexer(name string, input *bufio.Reader, start StateFunction) *Lexer {
	l := &Lexer{
		Name:      name,
		Input:     input,
		Buffer:    make([]rune, 0),
		Items:     make(chan Item),
		Current:   -1,
		readUntil: -1,
	}
	go l.Run(start)
	return l
}

// Run starts the lexing process. Tokens will be emitted on the Items channel.
func (l *Lexer) Run(start StateFunction) {
	for state := start; state != nil; {
		state = state(l)
	}
	close(l.Items) // No more tokens will be emitted.
}

// Emit emits a token to the waiting channel.
func (l *Lexer) Emit(t ItemType) {
	// Emit the token and update the position of the lexer against the input stream.
	length := l.Current - l.Start
	l.Items <- Item{
		Type:  t,
		Value: string(l.Buffer[0:length]),
	}
	l.Start = l.Current
	l.Buffer = l.Buffer[length:]
}

// Advance reads a rune from the Input and sets the current position.
func (l *Lexer) Advance() (rune, error) {
	// Check to see whether we already have it in the buffer, if so, read it from there.
	if l.Current+1 <= l.readUntil {
		l.Current++
		return l.Buffer[l.Current-l.Start], nil
	}

	r, _, err := l.Input.ReadRune()
	l.Buffer = append(l.Buffer, r)
	l.Current++
	l.readUntil = l.Current
	return r, err
}

// Retreat steps back a rune.
func (l *Lexer) Retreat() (rune, error) {
	newPos := l.Current - 1
	if newPos < -1 {
		l.Current = -1
		return 0x0, errors.New("cannot retreat past the start of the stream")
	}
	l.Current = newPos
	if l.Current > -1 {
		return l.Buffer[l.Current-l.Start], nil
	}
	return 0x0, nil
}

// StateFunction represents the state of the scanner as a function that returns
// the next state.
type StateFunction func(*Lexer) StateFunction
