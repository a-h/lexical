package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func any(pi parse.Input) parse.Result {
	_, err := pi.Advance()
	s := pi.Collect()
	return parse.String("any", s, any, err)
}

func uppercase(pi parse.Input) parse.Result {
	allowed := "ABCDEFGHIJKLMNOPQRSTUVXYZ"
	pr, err := pi.Peek()
	if strings.ContainsRune(allowed, pr) {
		_, err = pi.Advance()
		value := pi.Collect()
		return parse.String("uppercase", value, uppercase, err)
	}
	return parse.Failure()
}

func TestScanning(t *testing.T) {
	bs := bytes.NewBufferString("ABCDEFG")
	sr := bufio.NewReader(bs)
	stream := input.New("Scanner Input", sr)

	scanner := New(stream, any)
	var err error
	for {
		item, err := scanner.Next()
		if err != nil {
			break
		}
		fmt.Println(item)
		fmt.Println(item.Name())
	}
	if err != nil {
		t.Error(err)
	}
}
