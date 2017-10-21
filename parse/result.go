package parse

// Result is the result of a parse operation.
type Result interface {
	Success() bool
	Next() Function
	Item() ResultItem
	Error() error
}

// ResultItem is the value extracted by a parse operation.
type ResultItem interface {
	Name() string
	Value() interface{}
	Eq(ResultItem) bool
}

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
}

// Function represents the state of the scanner as a function that returns
// the next state.
type Function func(Input) Result
