package scanner

import (
	"fmt"
	"io"

	"github.com/a-h/lexical/parse"
)

// Scanner can take an input stream and execute parse results.
type Scanner struct {
	Input  parse.Input
	Parser parse.Function
}

// Next should be called repeatedly to request the next token from the stream.
// If
func (s *Scanner) Next() (item interface{}, err error) {
	result := s.Parser(s.Input)
	success := result.Success
	if !success && result.Error != io.EOF {
		line, col := s.Input.Position()
		return result.Item, fmt.Errorf("scanner: unmatched at line %v, column %v, item: %v", line, col, result)
	}
	s.Input.Collect()
	return result.Item, result.Error
}

// New creates a new Scanner.
func New(stream parse.Input, p parse.Function) *Scanner {
	return &Scanner{
		Input:  stream,
		Parser: p,
	}
}
