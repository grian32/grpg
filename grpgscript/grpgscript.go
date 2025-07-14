package main

import (
	"cmp"
	"io"
	"log"
	"os"
)

func RunFile(path string) {
	f, err1 := os.Open(path)
	bytes, err2 := io.ReadAll(f)

	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("Failed to run file with path %s %v", path, err)
	}

	Run(bytes)
}

func Run(bytes []byte) {
	// scanner := lex_old.NewScanner(string(bytes))
	// scanner.ScanTokens()
	// fmt.Println(lex_old.TokenSliceString(scanner.Tokens))
}
