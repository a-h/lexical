package scanner

import (
	"fmt"
	"io"

	"github.com/a-h/lexical/parse"
)

type Scanner struct {
	Input  parse.Input
	Parser parse.Function
}

func (s *Scanner) Next() parse.Result {
	result := s.Parser(s.Input)
	success := result.Success
	if !success && result.Error != io.EOF {
		line, col := s.Input.Position()
		result.Error = fmt.Errorf("scanner: unmatched at line %v, column %v", line, col)
	}
	if result.Error != io.EOF && result.Next == nil {
		line, col := s.Input.Position()
		result.Error = fmt.Errorf("scanner: not at end of file, but have finished parsing at line %v, column %v", line, col)
	}
	s.Parser = result.Next
	s.Input.Collect()
	return result
}

// New creates a new Scanner.
func New(stream parse.Input, p parse.Function) *Scanner {
	return &Scanner{
		Input:  stream,
		Parser: p,
	}
}
