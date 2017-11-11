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
		2.1,
		'一',
		"个",
	}
	result, _ := WithStringConcatCombiner(inputs)
	if result != "ABCDE12.1一个" {
		t.Errorf("Expected 'ABCDE12.1一个', but got '%v'", result)
	}
}

func BenchmarkWithStringConcatCombiner(b *testing.B) {
	items := []interface{}{'A', "BCDEF", 'G', "HIHJK"}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		WithStringConcatCombiner(items)
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
		WithIntegerCombiner([]interface{}{'1', "2", '3'})
	}
}
