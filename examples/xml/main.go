package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"unicode"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
	"github.com/a-h/lexical/scanner"
)

var profile = flag.Bool("profile", false, "Set to true to enable profiling to cpuprofile.out")
var memprofile = flag.Bool("memprofile", false, "Set to true to enable profiling to memprofile.out")

func main() {
	flag.Parse()

	filename := "example.xml"

	// Enable profiling.
	if *profile {
		f, err := os.Create("cpuprofile.out")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	count := countHouses(filename)

	if *memprofile {
		f, err := os.Create("memprofile.out")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	fmt.Printf("Go: Found %v houses\n", count)
}

func handle(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func countHouses(filename string) int {
	houses := 0

	file, err := os.Open(filename)
	handle(err)
	defer file.Close()
	buffer := bufio.NewReaderSize(file, 1024*1024*256) // 33554432
	scan := scanner.New(input.New(buffer), xmlTokens)
	for {
		item, err := scan.Next()
		switch v := item.(type) {
		case XMLSelfClosingElement:
			if v.Name == "ELEMENT" {
				houses++
			}
		default:
			fmt.Println(reflect.TypeOf(v))
		}
		if err != nil {
			break
		}
	}
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to parse file with err: %v\n", err)
	}

	return houses
}

type XMLWhitespace string

type XMLSelfClosingElement struct {
	Name       string
	Attributes []XMLAttribute
}

type XMLStartElement struct {
	Name       string
	Attributes []XMLAttribute
}

type XMLEndElement struct {
	Name string
}

type XMLAttribute struct {
	Name  string
	Value string
}

var xmlTokens = parse.Any(whiteSpace, selfClosingTag, startElement, closeElement, procInst)

var procInst = parse.String(`<?xml version="1.0" encoding="utf-8"?>`)

var asXMLWhitespace parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	v, ok := parse.WithStringConcatCombiner(inputs)
	s, _ := v.(string)
	return XMLWhitespace(s), ok
}

var whiteSpace = parse.AtLeast(asXMLWhitespace, 1, parse.RuneInRanges(unicode.White_Space))
var optionalWhiteSpace = parse.AtLeast(parse.WithStringConcatCombiner,
	0,
	parse.RuneInRanges(unicode.White_Space),
)
var letterOrDigit = parse.RuneInRanges(unicode.Letter, unicode.Number)

var xmlName = parse.Then(
	parse.WithStringConcatCombiner,
	parse.Letter,
	parse.Many(parse.WithStringConcatCombiner,
		0,   // minimum match count
		500, // maxmum match count
		letterOrDigit),
)

var asXMLSelfClosingElement parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	name, _ := inputs[1].(string)
	attributes, _ := inputs[2].([]XMLAttribute)
	return XMLSelfClosingElement{
		Name:       name,
		Attributes: attributes,
	}, true
}

var tagOpen = parse.Rune('<')
var tagOpenClosingTag = parse.String("</")
var tagClose = parse.Rune('>')
var tagSelfClose = parse.String("/>")

var selfClosingTag = parse.All(asXMLSelfClosingElement,
	tagOpen,
	xmlName, // 1: name
	parse.AtLeast(asXMLAttributeArray, 0,
		xmlAttribute,
	), // 2: attributes
	optionalWhiteSpace,
	tagSelfClose,
)

var asXMLAttribute parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	name, _ := inputs[1].(string)
	value, _ := inputs[6].(string)
	return XMLAttribute{
		Name:  name,
		Value: value,
	}, true
}

var equals = parse.Rune('=')
var quotes = parse.RuneIn(`"'`)

var xmlAttribute = parse.All(asXMLAttribute,
	whiteSpace,
	xmlName, // 1: name
	optionalWhiteSpace,
	equals,
	optionalWhiteSpace,
	quotes,
	parse.StringUntil(quotes), // 6: value
)

var asXMLAttributeArray parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	rv := make([]XMLAttribute, len(inputs))
	for i, v := range inputs {
		rv[i], _ = v.(XMLAttribute)
	}
	return rv, true
}

var asXMLStartElement parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	name, _ := inputs[1].(string)
	attributes, _ := inputs[2].([]XMLAttribute)
	return XMLStartElement{
		Name:       name,
		Attributes: attributes,
	}, true
}

var startElement = parse.All(asXMLStartElement,
	tagOpen,
	xmlName, // 1: name
	parse.AtLeast(asXMLAttributeArray, 0,
		xmlAttribute,
	), // 2: attributes
	tagClose,
)

var asXMLEndElement parse.MultipleResultCombiner = func(inputs []interface{}) (interface{}, bool) {
	name, _ := inputs[1].(string)
	return XMLEndElement{
		Name: name,
	}, true
}

var closeElement = parse.All(asXMLEndElement,
	tagOpenClosingTag,
	xmlName, // 1: name
	tagClose,
)
