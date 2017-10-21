package scanner

import (
	"fmt"

	"github.com/a-h/lexical/parse"
)

type Scanner struct {
	Input  parse.Input
	Parser parse.Function
}

func (s *Scanner) Next() parse.Result {
	result := s.Parser(s.Input)
	success := result.Success
	if !success && result.Error != nil {
		line, col := s.Input.Position()
		result.Error = fmt.Errorf("scanner: unmatched at line %v, column %v", line, col)
	}
	s.Parser = result.Next
	return result
}

// New creates a new Scanner.
func New(stream parse.Input, p parse.Function) *Scanner {
	return &Scanner{
		Input:  stream,
		Parser: p,
	}
}
