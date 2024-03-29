> **Warning**
> This library is now unmaintained in favour of a Go generics rewrite at https://github.com/a-h/parse

# Lexical

A set of parsing tools for Go inspired by [Sprache](https://github.com/sprache/Sprache/).

## Input

Parsers first need to read data to see if the data matches a pattern. If the data doesn't match, then the parser will need to move back to the end position of the last successful parse to try a different pattern.

Just storing everything in RAM works if your file sizes are small, but your process uses a lot of RAM.

Writing to a file to do this would mean seeking on disk, potentially making the performance suffer.

Instead, the `Stream` type provides a way of reading runes (characters) from an input `bufio.Reader` into a cache in RAM. Once a token has been consumed by the parser, the consumed bytes are discarded. The amount of RAM consumed will depend on the parser that uses it.

The `Stream` type implements the `parse.Input` interface:

```go
// Input represents the input to a parser.
type Input interface {
	// Collect collects all of the string data parsed so far and returns it, then starts a new collection
	// from the current position in the input.
	Collect() string
	// Advance advances the input by a single rune and consumes it.
	Advance() (rune, error)
	// Retreat retreats the input position by a single rune and unconsumes it.
	Retreat() (rune, error)
	// Peek returns the next rune from the input without consuming it.
	Peek() (rune, error)
	// Position returns the line and column number of the current position within the stream.
	Position() (line, column int)
	// Index returns the current index of the parser input.
	Index() int64
}
```

## Parser Functions

Parser functions provide a way of matching patterns in a given input. They are designed to be able to be composed together to make more complex operations.

The [examples](./examples) directory contains several examples of composing the primitive functions.

### Functions

* `Any`
    * Parse any of the provided parse functions, or roll back.
* `AnyRune`
    * Parse any rune.
* `AtLeast`
    * Parse the provided function at least the number of times specified, or roll back.
* `AtMost`
    * Parse the provided function at least once, and at most the number of times specified, or roll back.
* `Letter`
    * Parse any letter in the Unicode Letter range or roll back.
* `Many`
    * Parse the provided parse function a number of times or roll back.
* `Optional`
    * Attempt to parse, but don't roll back if a match isn't found.
* `Or`
    * Return the first successful result of the provided parse functions, or roll back.
* `Rune`
    * Parse the specified rune (character) or fallback.
* `RuneIn`
    * Parse a rune from the input stream if it's in the specified string, or roll back.
* `RuneInRanges`
    * Parse a rune from the input stream if it's in the specified Unicode ranges, or roll back.
* `RuneNotIn`
    * Parse a rune from the input stream if it's not in the specified string, or roll back.
* `RuneWhere`
    * Parse a rune from the input stream if the predicate function passed in succeeds, or roll back.
* `String`
    * Parse a string from the input stream if it exactly matches the provided string, or roll back.
* `StringUntil`
    * Parse a string from the input stream until the specified _until_ parser is matched.
* `Then`
    * Return the results of the first and second parser passed through the combiner function which converts the two results into a single output (a map / reduce operation), or roll back if either doesn't match.
* `Times`
    * Parse using the specified function a set number of times or roll back.
* `ZeroToNine`
    * Parse a rune from the input stream if it's within the set of 1234567890.

### Examples

Using the `Or` function to parse either 'A' or 'B':

```go
parser := parse.Or(parse.Rune('A'), parse.Rune('B'))

matchesA := parser(input.NewFromString("A")).Success // true
matchesB := parser(input.NewFromString("B")).Success // true
matchesC := parser(input.NewFromString("C")).Success // false

fmt.Println(matchesA) // true
fmt.Println(matchesB) // true
fmt.Println(matchesC) // false

```

The `Or` function only returns a single result but the `Many` function is more complex, because you generally want to do something with the results, such as convert the runes or strings captured by the parser into another value. The `parse.WithIntegerCombiner` and `parse.WithStringConcatCombiner` functions provide some default implementations.

The [examples](./examples) directory contains several examples of taking the primitive parse results and returning other types such as dates and URLs.


```go
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

// Parse letters into a string
upToThreeLetters := parse.AtMost(parse.WithStringConcatCombiner, 3, parse.Letter)
letters := upToThreeLetters(input.NewFromString("ABC1"))
// Check that we got a string back from the parser with a 'type assertion'.
resultItem, ok := letters.Item.(string)
if !ok || resultItem != "ABC" {
    t.Errorf("for 'ABC1', expected to extract 'ABC', but extracted '%v'", letters.Item)
}
```

## Scanner

The `Scanner` type combines the parser functions and `Stream` type to allow parsing of input files. See `scanner_test.go` for a working example.

```go
stream := input.NewFromString(`<a>Example</a>`)

scanner := New(stream, xmlTokens)
var err error
for {
    item, err := scanner.Next()
    // Do something with the results based on the 
    // token's type.
    switch v := item.(type) {
        case string:
            fmt.Println(v)
        case int:
            fmt.Println(v)
    }
    if err != nil {
        break
    }
}
if err != nil && err != io.EOF {
    panic("error")
}
```
