package examples

import (
	"testing"
	"time"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestDateParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedMatch bool
		expectedValue time.Time
	}{
		{
			input:         "2001-02-03",
			expectedMatch: true,
			expectedValue: time.Date(2001, 02, 03, 0, 0, 0, 0, time.UTC),
		},
		{
			input:         "clearly not a date",
			expectedMatch: false,
		},
		{
			input:         "201-01-1", // Missing some digits
			expectedMatch: false,
		},
	}

	for _, test := range tests {
		ip := input.NewFromString("date_test", test.input)
		result := date(ip)
		if result.Success != test.expectedMatch {
			t.Errorf("for input '%v', expected sucess %v, but was %v", test.input, test.expectedMatch, result.Success)
		}
	}
}

var year = parse.Many(parse.WithIntegerCombiner, 4, 4, parse.ZeroToNine)
var month = parse.All(parse.WithIntegerCombiner, parse.RuneIn("01"), parse.ZeroToNine)
var day = parse.All(parse.WithIntegerCombiner, parse.RuneIn("0123"), parse.ZeroToNine)

var date = parse.All(dateConverter, year, parse.Rune('-'), month, parse.Rune('-'), day)

func dateConverter(items []interface{}) (item interface{}, value bool) {
	yi, mi, di := items[0], items[2], items[4]
	y, ok := yi.(int)
	if !ok {
		return 0, false
	}
	m, ok := mi.(int)
	if !ok {
		return 0, false
	}
	d, ok := di.(int)
	if !ok {
		return 0, false
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), true
}
