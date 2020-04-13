package input

import "testing"

func TestFromBufferFunction(t *testing.T) {
	tests := []struct {
		name               string
		startOfBufferIndex int64
		currentIndex       int64
		buffer             *Buffer
		expectedRune       rune
		expectedOK         bool
	}{
		{
			name:               "Read A from 'ABC'",
			startOfBufferIndex: 0,
			currentIndex:       1,
			buffer:             NewBufferFromString("ABC", 3),
			expectedRune:       'A',
			expectedOK:         true,
		},
		{
			name:               "Read B from 'ABC'",
			startOfBufferIndex: 0,
			currentIndex:       2,
			buffer:             NewBufferFromString("ABC", 3),
			expectedRune:       'B',
			expectedOK:         true,
		},
		{
			name:               "Read C from 'ABC'",
			startOfBufferIndex: 0,
			currentIndex:       3,
			buffer:             NewBufferFromString("ABC", 3),
			expectedRune:       'C',
			expectedOK:         true,
		},
		{
			name:               "Read D from 'ABC'",
			startOfBufferIndex: 0,
			currentIndex:       4,
			buffer:             NewBufferFromString("ABC", 3),
			expectedRune:       0x0,
			expectedOK:         false,
		},
	}

	for _, test := range tests {
		r, ok := fromBuffer(test.startOfBufferIndex, test.currentIndex, test.buffer)
		if r != test.expectedRune {
			t.Errorf("%s: expected rune '%v' but got '%v'", test.name, string(test.expectedRune), string(r))
		}
		if ok != test.expectedOK {
			t.Errorf("%s: expected to read from buffer to be %v, but was %v", test.name, test.expectedOK, ok)
		}
	}
}
