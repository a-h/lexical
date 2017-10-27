package parse

import "fmt"

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

// Function represents the state of the scanner as a function that returns
// the next state.
type Function func(Input) Result

// Result is the result of a parse operation.
type Result struct {
	Name    string
	Success bool
	Item    interface{}
	Error   error
}

// Success creates a successful result of a parse operation.
func Success(name string, item interface{}, err error) Result {
	return Result{
		Name:    name,
		Success: true,
		Item:    item,
		Error:   err,
	}
}

// Failure creates an unsuccessful result of a parse operation.
func Failure(name string, err error) Result {
	return Result{
		Name:    name,
		Success: false,
		Error:   err,
	}
}

// Eq compares two results for equality.
func (result Result) Eq(cmp Result) bool {
	if cmp.Name != result.Name {
		return false
	}
	if cmp.Item != result.Item {
		return false
	}
	return true
}

// String returns the string representation of a result, truncated to 10 characters.
func (result Result) String() string {
	if !result.Success {
		return fmt.Sprintf("✗ (%v) err: %v", result.Name, result.Error)
	}

	success := "✗"
	if result.Success {
		success = "✓"
	}

	var v string
	switch ti := result.Item.(type) {
	case rune:
		v = string(ti)
	case string:
		v = ti
	default:
		v = fmt.Sprintf("%v", ti)
	}

	if len(v) > 13 {
		v = v[0:10] + "..."
	}

	e := ""
	if result.Error != nil {
		e = fmt.Sprintf("\n%v (%v) err: %v", success, result.Name, result.Error)
	}

	return fmt.Sprintf("%v (%v) %v%v", success, result.Name, v, e)
}
