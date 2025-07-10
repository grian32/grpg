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
)

func init() {
	// TODO: surely i can abstract this somehow :S
	file, err1 := os.Open("../testdata/symbols.grpgscript")
	bytes, err2 := io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("Error reading symbols file: %v", err)
	}
	singleSymbols = bytes

	file, err1 = os.Open("../testdata/doublesymbols.grpgscript")
	bytes, err2 = io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("Error reading double symbols file: %v", err)
	}
	doubleSymbols = bytes

	file, err1 = os.Open("../testdata/symbolscomments.grpgscript")
	bytes, err2 = io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("Error reading symbols comments file: %v", err)
	}
	symbolComments = bytes
}

func BenchmarkSingleSymbols(b *testing.B) {
	scanner := NewScanner(string(singleSymbols))
	scanner.ScanTokens()
}

func BenchmarkDoubleSymbols(b *testing.B) {
	scanner := NewScanner(string(doubleSymbols))
	scanner.ScanTokens()
}

func BenchmarkSymbolsComments(b *testing.B) {
	scanner := NewScanner(string(symbolComments))
	scanner.ScanTokens()
}
