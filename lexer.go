package lexical

import (
	"bufio"
	"errors"
	"fmt"
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
	// CurrentRune is the current rune at the cursor.
	CurrentRune rune
}

func (l *Lexer) String() string {
	return fmt.Sprintf("%v: Current Rune: '%v', Start of Token Position: %v, Current Position: %v, Forward Buffer Size: %v, Current Buffer: '%v'", l.Name, string(l.CurrentRune), l.Start, l.Current, len(l.Buffer), string(l.Buffer))
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
	lastPosition := l.Current
	for state := start; state != nil; {
		state = state(l)
		if l.Current == lastPosition {
			l.EmitError(fmt.Errorf("lexer: stuck in a loop at position %v", l.Current))
			return
		}
		lastPosition = l.Current
	}
	close(l.Items) // No more tokens will be emitted.
}

// Emit emits a token to the waiting channel.
func (l *Lexer) Emit(t ItemType) Item {
	// Emit the token and update the position of the lexer against the input stream.
	// Returning the item helps with unit testing.
	length := int(l.Current - l.Start)
	left := getLeft(l.Buffer, length)
	right := getRight(l.Buffer, length)
	item := Item{
		Type:  t,
		Value: string(left),
	}
	l.Items <- item
	l.Start = l.Current
	l.Buffer = right
	return item
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

// EmitError emits an error to the waiting channel.
func (l *Lexer) EmitError(err error) {
	l.Items <- Item{
		Type:  ItemTypeError,
		Value: err.Error(),
	}
	close(l.Items)
}

// Advance reads a rune from the Input and sets the current position.
func (l *Lexer) Advance() (rune, error) {
	// Check to see whether we already have it in the buffer, if so, read it from there.
	if l.Current+1 <= l.readUntil {
		l.Current++
		l.CurrentRune = l.Buffer[l.Current-l.Start]
		return l.CurrentRune, nil
	}

	r, _, err := l.Input.ReadRune()
	l.Buffer = append(l.Buffer, r)
	l.Current++
	l.readUntil = l.Current
	l.CurrentRune = r
	return r, err
}

// Peek reads a rune from the Input, then sets the current position back.
func (l *Lexer) Peek() (rune, error) {
	r, err := l.Advance()
	if err != nil {
		return r, fmt.Errorf("lexer.peek: failed to advance: %v", err)
	}
	_, err = l.Retreat()
	return r, err
}

// Retreat steps back a rune.
func (l *Lexer) Retreat() (rune, error) {
	newPos := l.Current - 1
	if newPos < -1 {
		l.Current = -1
		l.CurrentRune = 0x0
		return 0x0, errors.New("cannot retreat past the start of the stream")
	}
	l.Current = newPos
	if l.Current > -1 {
		l.CurrentRune = l.Buffer[l.Current-l.Start]
		return l.CurrentRune, nil
	}
	return 0x0, nil
}

// StateFunction represents the state of the scanner as a function that returns
// the next state.
type StateFunction func(*Lexer) StateFunction

// AdvanceUntilRune advances the reader until the rune r is reached.
func (l *Lexer) AdvanceUntilRune(r rune) (err error) {
	return l.AdvanceUntil(func(rr rune) bool { return r == rr })
}

// AdvanceUntil advances the reader until the rule function returns true.
func (l *Lexer) AdvanceUntil(rule func(r rune) bool) (err error) {
	for {
		cr, err := l.Advance()
		if err != nil {
			return err
		}
		if rule(cr) {
			return nil
		}
	}
}
