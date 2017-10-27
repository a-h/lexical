package examples

import (
	"fmt"
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestCompose(t *testing.T) {
	parser := parse.Or(parse.Rune('A'), parse.Rune('B'))

	matchesA := parser(input.NewFromString("A")).Success // true
	matchesB := parser(input.NewFromString("B")).Success // true
	matchesC := parser(input.NewFromString("C")).Success // false

	fmt.Println(matchesA) // true
	fmt.Println(matchesB) // true
	fmt.Println(matchesC) // false
}

func TestMany(t *testing.T) {
	// parse.WithIntegerCombiner concatentates the captured runes into a string,
	// and parses the result to an integer.
	oneToThreeNumbers := parse.Many(parse.WithIntegerCombiner,
		1, // minimum match count
		3, // maximum match count
		parse.ZeroToNine)

	resultA := oneToThreeNumbers(input.NewFromString("123"))
	fmt.Println(resultA.Success) // true
	fmt.Println(resultA.Item)    // integer value of 123

	resultB := oneToThreeNumbers(input.NewFromString("1234"))
	fmt.Println(resultB.Success) // true
	fmt.Println(resultB.Item)    // integer value of 123

	// This Many function will stop reading at the 'a'.
	resultC := oneToThreeNumbers(input.NewFromString("1a234"))
	fmt.Println(resultC.Success) // true
	fmt.Println(resultC.Item)    // integer value of 1
}
