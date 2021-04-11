package scanner

import (
	"io"
	"reflect"
	"testing"
	"unicode"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestScanningAnyRunes(t *testing.T) {
	text := `abcdef`
	stream := input.NewFromString(text)

	var output string
	scanner := New(stream, parse.AnyRune())
	var i int
	var err error
	for {
		item, err := scanner.Next()
		if err != nil {
			break
		}
		output += string(item.(rune))
		i++
		if i > 10 {
			t.Errorf("infinite loop?")
			break
		}
	}
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if i != 6 {
		t.Errorf("expected to read 6 runes, got %d", i)
	}
	if output != text {
		t.Errorf("expected %q, got %q", text, output)
	}
}

func TestScanning(t *testing.T) {
	text := `<a>abc</a><b>def</b>`
	stream := input.NewFromString(text)

	scanner := New(stream, xmlTag)
	var i int
	var err error
	for {
		_, err := scanner.Next()
		if err != nil {
			break
		}
		i++
		if i > 10 {
			t.Errorf("infinite loop")
			break
		}
	}
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if i != 2 {
		t.Errorf("expected to read 2 tags (2 x open, text, and close), but got %d", i)
	}
}

var xmlTag = parse.All(parse.WithStringConcatCombiner, xmlOpenElement, xmlText, xmlCloseElement)

var combineTagAndContents parse.MultipleResultCombiner = func(results []interface{}) (interface{}, bool) {
	name, ok := results[1].(string)
	return "name: " + name, ok
}

var letterOrDigit = parse.RuneInRanges(unicode.Letter, unicode.Number)

var xmlName = parse.Then(
	parse.WithStringConcatCombiner,
	parse.RuneWhere(unicode.IsLetter),
	parse.Many(parse.WithStringConcatCombiner,
		0,   // minimum match count
		500, // maximum match count
		letterOrDigit),
)

var xmlOpenElement = parse.All(
	combineTagAndContents,
	parse.Rune('<'),
	parse.StringUntil(parse.Rune('>')),
	parse.Rune('>'),
)

var xmlText = parse.StringUntil(parse.String("<"))

var xmlCloseElement = parse.All(
	parse.WithStringConcatCombiner,
	parse.String("</"),
	parse.StringUntil(parse.String(">")),
	parse.String(">"),
)

func TestXMLName(t *testing.T) {
	tests := []struct {
		input        string
		expected     bool
		expectedItem string
	}{
		{
			input:        "AB",
			expected:     true,
			expectedItem: "AB",
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		result := xmlName(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if test.expected && result.Item != test.expectedItem {
			t.Errorf("test %v: for input '%v' expected item '%v' but got '%v' (%v)", i, test.input, test.expectedItem, result.Item, reflect.TypeOf(result.Item))
		}
	}
}
