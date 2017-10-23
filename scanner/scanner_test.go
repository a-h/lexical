package scanner

import (
	"fmt"
	"strings"
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func any(pi parse.Input) parse.Result {
	_, err := pi.Advance()
	s := pi.Collect()
	return parse.Success("any", s, any, err)
}

func uppercase(pi parse.Input) parse.Result {
	allowed := "ABCDEFGHIJKLMNOPQRSTUVXYZ"
	pr, err := pi.Peek()
	if strings.ContainsRune(allowed, pr) {
		_, err = pi.Advance()
		value := pi.Collect()
		return parse.Success("uppercase", value, uppercase, err)
	}
	return parse.Failure("uppercase", fmt.Errorf("Expected A-Z, but got '%v'", pr))
}

func TestScanning(t *testing.T) {
	stream := input.NewFromString("Scanner Input", "ABCDEFG")

	scanner := New(stream, any)
	var err error
	for {
		result := scanner.Next()
		if result.Error != nil {
			break
		}
		fmt.Println(result)
	}
	if err != nil {
		t.Error(err)
	}
}
