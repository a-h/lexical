package input

import "fmt"

// Position represents the character position within a text file.
type Position struct {
	Index int64
	Line  int
	Col   int
	// The lengths of each line we've seen.
	lineLengths     map[int]int
	carriageReturns map[int64]struct{}
	lineFeeds       map[int64]struct{}
}

// NewPosition creates a Position to represent the character position within a text file.
func NewPosition(line int, col int) Position {
	return Position{
		Index:           int64(-1),
		Line:            line,
		Col:             col,
		lineLengths:     make(map[int]int),
		carriageReturns: make(map[int64]struct{}),
		lineFeeds:       make(map[int64]struct{}),
	}
}

// String creates a string representation of the position.
func (p *Position) String() string {
	return fmt.Sprintf("Line: %v, Col: %v", p.Line, p.Col)
}

// Eq compares two positions and returns true if they are equal.
func (p *Position) Eq(cmp Position) bool {
	return p.Col == cmp.Col && p.Line == cmp.Line
}

// Advance advances the position by a line if the rune is'\n', does nothing if the rune
// is '\r' and advances by a col character if the rune is anything else.
func (p *Position) Advance(r rune) {
	p.Index++
	if r == '\r' {
		p.carriageReturns[p.Index] = struct{}{}
		return
	}
	if r == '\n' {
		p.lineLengths[p.Line] = p.Col
		p.lineFeeds[p.Index] = struct{}{}
		p.Line++
		p.Col = 0
		return
	}
	p.Col++
}

// Retreat decreases the position by a line if the rune is'\n', does nothing if the rune
// is '\r' and decreases by a col character if the rune is anything else.
func (p *Position) Retreat(r rune) {
	p.Index--
	if r == '\r' {
		return
	}
	lfIndex := p.Index + 1
	if _, isRetreatingFromCR := p.carriageReturns[p.Index+1]; isRetreatingFromCR {
		lfIndex++
	}
	if _, isRetreatingFromNewLine := p.lineFeeds[lfIndex]; isRetreatingFromNewLine {
		p.Line--
		p.Col = p.lineLengths[p.Line]
		return
	}
	p.Col--
}
