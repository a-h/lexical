package scanner

import (
	"fmt"
	"io"
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestScanning(t *testing.T) {
	stream := input.NewFromString("Scanner Input", `<a>Example</a>`)

	scanner := New(stream, xmlTokens)
	var err error
	for {
		result := scanner.Next()
		if result.Error != nil {
			err = result.Error
			break
		}
		fmt.Printf("Result: %v\n", result)
		stream.Collect()
	}
	if err != nil && err != io.EOF {
		t.Error(err)
	}
}

var xmlTokens = parse.Any(xmlOpenElement, xmlText, xmlCloseElement)

var combineTagAndContents parse.ResultCombiner = func(openTag interface{}, tagContents interface{}) interface{} {
	return "name: " + fmt.Sprintf("%v", tagContents)
}

var xmlOpenElement = parse.Then(
	parse.Rune('<'),
	parse.StringUntil(parse.Rune('>')),
	combineTagAndContents,
)

var xmlText = parse.StringUntil(parse.Rune('<'))

var xmlCloseElement = parse.All(
	parse.ConcatenateStringsCombiner,
	parse.Rune('<'),
	parse.Rune('/'),
	parse.StringUntil(parse.Rune('>')),
	parse.Rune('>'),
)
