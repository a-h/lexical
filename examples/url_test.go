package examples

import (
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestURLParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedMatch bool
		expectedValue URL
	}{
		{
			input:         "example.com",
			expectedMatch: true,
			expectedValue: URL{
				Host: "example.com",
			},
		},
		{
			input:         "example.com:80/path",
			expectedMatch: true,
			expectedValue: URL{
				Host: "example.com:80",
				Path: "/path",
			},
		},
	}

	for _, test := range tests {
		ip := input.NewFromString("url_test", test.input)
		result := url(ip)
		if result.Success != test.expectedMatch {
			t.Errorf("for input '%v', expected success '%v', but was %v", test.input, test.expectedMatch, result.Success)
		}
		if result.Item != test.expectedValue {
			t.Errorf("for input '%v', expected sucess '%v', but was %v", test.input, test.expectedValue, result.Item)
		}
	}
}

type URL struct {
	Host string
	Path string
}

var head = parse.Any(parse.Letter, parse.ZeroToNine)
var tailCharacter = parse.Any(parse.Letter, parse.ZeroToNine, parse.Rune('-'))
var tail = parse.Many(parse.WithStringConcatCombiner, 0, 63, tailCharacter)

var host = parse.All(parse.WithStringConcatCombiner, head, tail)

var singleTLD = parse.All(parse.WithStringConcatCombiner,
	parse.Rune('.'),
	parse.Many(parse.WithStringConcatCombiner, 0, 3, head),
)

var allTLDs = parse.Many(parse.WithStringConcatCombiner, 1, 2, singleTLD)

var port = parse.Many(parse.WithStringConcatCombiner, 0, 1,
	parse.All(parse.WithStringConcatCombiner,
		parse.Rune(':'),
		parse.Many(parse.WithStringConcatCombiner, 1, 5, parse.ZeroToNine),
	),
)

var domain = parse.All(parse.WithStringConcatCombiner,
	host,
	allTLDs,
	port,
)

var path = parse.All(parse.WithStringConcatCombiner,
	parse.Rune('/'),
	parse.Many(parse.WithStringConcatCombiner, 0, 65536, parse.AnyRune()),
)

var optionalPath = parse.Optional(parse.WithStringConcatCombiner, path)

func urlParser(results []interface{}) (rv interface{}, ok bool) {
	host, _ := results[0].(string)
	path, _ := results[1].(string)

	ok = true
	url := URL{
		Host: host,
		Path: path,
	}
	return url, ok
}

var url = parse.All(urlParser, domain, optionalPath)
