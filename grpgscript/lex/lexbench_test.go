package lex

import (
	"cmp"
	"io"
	"log"
	"os"
	"testing"
)

var (
	singleSymbols  []byte
	doubleSymbols  []byte
	symbolComments []byte
	numbers        []byte
	strings        []byte
	helloWorld     []byte
	functions      []byte
)

func init() {
	files := []string{
		"../testdata/symbols.grpgscript",
		"../testdata/doublesymbols.grpgscript",
		"../testdata/symbolscomments.grpgscript",
		"../testdata/numbers.grpgscript",
		"../testdata/strings.grpgscript",
		"../testdata/helloworld.grpgscript",
		"../testdata/functions.grpgscript",
	}
	// doesn't work if i store this in a struct, idk why :(
	arrays := []*[]byte{
		&singleSymbols,
		&doubleSymbols,
		&symbolComments,
		&numbers,
		&strings,
		&helloWorld,
		&functions,
	}

	for idx, name := range files {
		file, err1 := os.Open(name)
		bytes, err2 := io.ReadAll(file)
		if err := cmp.Or(err1, err2); err != nil {
			log.Fatalf("Error reading %s file: %v", name, err)
		}
		// required to modify the arr
		*arrays[idx] = bytes
	}
}

func BenchmarkLex(b *testing.B) {
	benches := []struct {
		name string
		data string
	}{
		{"SingleSymbols", string(singleSymbols)},
		{"DoubleSymbols", string(doubleSymbols)},
		{"SymbolsComments", string(symbolComments)},
		{"Numbers", string(numbers)},
		{"Strings", string(strings)},
		{"HelloWorld", string(helloWorld)},
		{"Functions", string(functions)},
	}

	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			scanner := NewScanner(bench.data)
			scanner.ScanTokens()
		})
	}
}
