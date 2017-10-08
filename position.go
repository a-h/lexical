package lexical

import "fmt"

// Position represents the character position within a text file.
type Position struct {
	Line int
	Col  int
	// The lengths of each line we've seen.
	lineLengths map[int]int
}

// NewPosition creates a Position to represent the character position within a text file.
func NewPosition(line int, col int) Position {
	return Position{
		Line:        line,
		Col:         col,
		lineLengths: make(map[int]int),
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
	if r == '\r' {
		return
	}
	if r == '\n' {
		// Store the line length for when we retreat.
		p.lineLengths[p.Line] = p.Col
		p.Line++
		p.Col = 0
		return
	}
	p.Col++
}

// Retreat decreases the position by a line if the rune is'\n', does nothing if the rune
// is '\r' and decreases by a col character if the rune is anything else.
func (p *Position) Retreat(r rune) {
	if r == '\r' {
		return
	}
	if r == '\n' {
		p.Line--
		// Retrieve the line length.
		p.Col = p.lineLengths[p.Line]
		return
	}
	p.Col--
}
