package parse

import (
	"testing"
)

func TestWithStringConcatCombiner(t *testing.T) {
	inputs := []interface{}{
		'A',
		"BCD",
		'E',
		1,
		2.0,
	}
	result, _ := WithStringConcatCombiner(inputs)
	if result != "ABCDE12" {
		t.Errorf("Expected 'ABCDE12.0', but got '%v'", result)
	}
}

func BenchmarkWithStringConcatCombiner(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		WithStringConcatCombiner([]interface{}{'A', "BCDEF", 'G', "HIHJK"})
	}
}

func TestWithIntegerCombiner(t *testing.T) {
	inputs := []interface{}{
		'1',
		'2',
		"30",
	}
	result, _ := WithIntegerCombiner(inputs)
	if result != 1230 {
		t.Errorf("Expected 1230, but got '%v'", result)
	}
}

func BenchmarkWithIntegerCombiner(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		WithIntegerCombiner([]interface{}{1, 2, 3})
	}
}
