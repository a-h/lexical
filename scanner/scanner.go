package scanner

import (
	"fmt"

	"github.com/a-h/lexical/parse"
)

type Scanner struct {
	Input  parse.Input
	Parser parse.Function
}

func (s *Scanner) Next() (parse.ResultItem, error) {
	result := s.Parser(s.Input)
	success := result.Success()
	if !success {
		line, col := s.Input.Position()
		return nil, fmt.Errorf("scanner: unmatched at line %v, column %v", line, col)
	}
	s.Parser = result.Next()
	return result.Item(), result.Error()
}

// New creates a new Scanner.
func New(stream parse.Input, p parse.Function) *Scanner {
	return &Scanner{
		Input:  stream,
		Parser: p,
	}
}
